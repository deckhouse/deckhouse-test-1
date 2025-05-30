---
title: "K/V Version 2"
permalink: en/stronghold/documentation/user/secrets-engines/kv/kv-v2.html
lang: en
description: The KV secrets engine can store arbitrary secrets.
---

The `kv` secrets engine is used to store arbitrary secrets within the
configured physical storage for Stronghold.

Key names must always be strings. If you write non-string values directly via
the CLI, they will be converted into strings. However, you can preserve
non-string values by writing the key/value pairs to Stronghold from a JSON file or
using the HTTP API.

This secrets engine honors the distinction between the `create` and `update`
capabilities inside ACL policies. The `patch` capability is also supported
which is used to represent partial updates whereas the `update` capability
represents full overwrites.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

A v2 `kv` secrets engine can be enabled by:

```shell-session
d8 stronghold secrets enable -version=2 kv
```

Or, you can pass `kv-v2` as the secrets engine type:

```shell-session
d8 stronghold secrets enable kv-v2
```

Additionally, when running a dev-mode server, the v2 `kv` secrets engine is enabled by default at the
path `secret/` (for non-dev servers, it is currently v1). It can be disabled, moved, or enabled multiple times at
different paths. Each instance of the KV secrets engine is isolated and unique.

## Upgrading from version 1

An existing version 1 kv store can be upgraded to a version 2 kv store via the CLI or API, as shown below. This will start an upgrade process to upgrade the existing key/value data to a versioned format. The mount will be inaccessible during this process. This process could take a long time, so plan accordingly.

Once upgraded to version 2, the former paths at which the data was accessible will no longer suffice. You will need to adjust user policies to add access to the version 2 paths as detailed in the [ACL Rules section below](/docs/secrets/kv/kv-v2#acl-rules). Similarly, users/applications will need to update the paths at which they interact with the kv data once it has been upgraded to version 2.

An existing version 1 kv can be upgraded to a version 2 KV store with the CLI command:

```shell-session
$ d8 stronghold kv enable-versioning secret/
Success! Tuned the secrets engine at: secret/
```

or via the API:

```shell-session
$ cat payload.json
{
  "options": {
      "version": "2"
  }
}
```

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/mounts/secret/tune
```

## ACL rules

The version 2 kv store uses a prefixed API, which is different from the
version 1 API. Before upgrading from a version 1 kv the ACL rules
should be changed. Also different paths in the version 2 API can be ACL'ed
differently.

Writing and reading versions are prefixed with the `data/` path. This policy
that worked for the version 1 kv:

```plaintext
path "secret/dev/team-1/*" {
  capabilities = ["create", "update", "read"]
}
```

Should be changed to:

```plaintext
path "secret/data/dev/team-1/*" {
  capabilities = ["create", "update", "read"]
}
```

There are different levels of data deletion for this backend. To grant a policy
the permissions to delete the latest version of a key:

```plaintext
path "secret/data/dev/team-1/*" {
  capabilities = ["delete"]
}
```

To allow the policy to delete any version of a key:

```plaintext
path "secret/delete/dev/team-1/*" {
  capabilities = ["update"]
}
```

To allow a policy to undelete data:

```plaintext
path "secret/undelete/dev/team-1/*" {
  capabilities = ["update"]
}
```

To allow a policy to destroy versions:

```plaintext
path "secret/destroy/dev/team-1/*" {
  capabilities = ["update"]
}
```

To allow a policy to list keys:

```plaintext
path "secret/metadata/dev/team-1/*" {
  capabilities = ["list"]
}
```

To allow a policy to view metadata for each version:

```plaintext
path "secret/metadata/dev/team-1/*" {
  capabilities = ["read"]
}
```

To allow a policy to permanently remove all versions and metadata for a key:

```plaintext
path "secret/metadata/dev/team-1/*" {
  capabilities = ["delete"]
}
```

The `allowed_parameters`, `denied_parameters`, and `required_parameters` fields are
not supported for policies used with the version 2 kv store. See the [Policies Concepts](/docs/concepts/policies)
for a description of these parameters.

See the [API Specification](/api-docs/secret/kv/kv-v2) for more
information.

## Usage

After the secrets engine is configured and a user/machine has an Stronghold token with
the proper permission, it can generate credentials. The `kv` secrets engine
allows for writing keys with arbitrary values.

The path-like KV-v1 syntax for referencing a secret (`secret/foo`) can still
be used in KV-v2, but we recommend using the `-mount=secret` flag syntax to
avoid mistaking it for the actual path to the secret (`secret/data/foo` is the
real path).

### Writing/Reading arbitrary data

1. Write arbitrary data:

   ```shell-session
   $ d8 stronghold kv put -mount=secret my-secret foo=a bar=b
   Key              Value
   ---              -----
   created_time     2024-06-19T17:20:22.985303Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          1
   ```

1. Read arbitrary data:

   ```shell-session
   $ d8 stronghold kv get -mount=secret my-secret
   ====== Metadata ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:20:22.985303Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          1

   ====== Data ======
   Key         Value
   ---         -----
   foo         a
   bar         b
   ```

1. Write another version, the previous version will still be accessible. The
   `-cas` flag can optionally be passed to perform a check-and-set operation. If
   not set the write will be allowed. In order for a write to be successful, `cas` must be set to
  the current version of the secret. If set to 0 a write will only be allowed if
  the key doesn’t exist as unset keys do not have any version information. Also
  remember that soft deletes do not remove any underlying version data from storage.
  In order to write to a soft deleted key, the cas parameter must match the key's
  current version.

   ```shell-session
   $ d8 stronghold kv put -mount=secret -cas=1 my-secret foo=aa bar=bb
   Key              Value
   ---              -----
   created_time     2024-06-19T17:22:23.369372Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          2
   ```

1. Reading now will return the newest version of the data:

   ```shell-session
   $ d8 stronghold kv get -mount=secret my-secret
   ====== Metadata ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:22:23.369372Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          2

   ====== Data ======
   Key         Value
   ---         -----
   foo         aa
   bar         bb
   ```

1. Partial updates can be accomplished using the `d8 stronghold kv patch` command. A
   command will initially attempt an HTTP `PATCH` request which requires the
   `patch` ACL capability. The `PATCH` request will fail if the token used
   is associated with a policy that does not contain the `patch` capability. In
   this case the command will perform a read, local update, and subsequent write
   which require both the `read` and `update` ACL capabilities.

   The `-cas` flag can optionally be passed to perform a check-and-set operation.
   It will only be used in the case of the initial `PATCH` request. The
   read-then-write flow will use the `version` value from the secret returned by
   the read to perform a check-and-set operation in the subsequent write.

   ```shell-session
   $ d8 stronghold kv patch -mount=secret -cas=2 my-secret bar=bbb
   Key              Value
   ---              -----
   created_time     2024-06-19T17:23:49.199802Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          3
   ```

1. The `d8 stronghold kv patch` command also supports a `-method` flag which can be
   used to specify HTTP `PATCH` or read-then-write. The supported values are
   `patch` and `rw` for HTTP `PATCH` and read-then-write, respectively.

   Perform a patch using the `patch` method:

   ```shell-session
   $ d8 stronghold kv patch -mount=secret -method=patch -cas=2 my-secret bar=bbb
   Key              Value
   ---              -----
   created_time     2024-06-19T17:23:49.199802Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          3
   ```

   Perform a patch using the read-then-write method:

   ```shell-session
   $ d8 stronghold kv patch -mount=secret -method=rw my-secret bar=bbb
   Key              Value
   ---              -----
   created_time     2024-06-19T17:23:49.199802Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          3
   ```

1. Reading after a patch will return the newest version of the data in which
   only the specified fields were updated:

   ```shell-session
   $ d8 stronghold kv get -mount=secret my-secret
   ====== Metadata ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:23:49.199802Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          3

   ====== Data ======
   Key         Value
   ---         -----
   foo         aa
   bar         bbb
   ```

1. Previous versions can be accessed with the `-version` flag:

   ```shell-session
   $ d8 stronghold kv get -mount=secret -version=1 my-secret
   ====== Metadata ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:20:22.985303Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          1

   ====== Data ======
   Key         Value
   ---         -----
   foo         a
   bar         b
   ```

You can also use [Stronghold's password policy](/docs/concepts/password-policies) feature to generate arbitrary values.

1. Write a password policy:

   ```shell-session
   $ d8 stronghold write sys/policies/password/example policy=-<<EOF

     length=20

     rule "charset" {
       charset = "abcdefghij0123456789"
       min-chars = 1
     }

     rule "charset" {
       charset = "!@#$%^&*STUVWXYZ"
       min-chars = 1
     }

   EOF
   ```

1. Write data using the `example` policy:

   ```shell-session
   $ d8 stronghold kv put -mount=secret my-generated-secret \
       password=$(d8 stronghold read -field password sys/policies/password/example/generate)
   ```

   **Example output:**

   ```plaintext
   ========= Secret Path =========
   secret/data/my-generated-secret

   ======= Metadata =======
   Key                Value
   ---                -----
   created_time       2023-05-10T14:32:32.37354939Z
   custom_metadata    <nil>
   deletion_time      n/a
   destroyed          false
   version            1
   ```

1. Read the generated data:

   ```shell-session
   $ d8 stronghold kv get -mount=secret my-generated-secret
   ========= Secret Path =========
   secret/data/my-generated-secret

   ======= Metadata =======
   Key                Value
   ---                -----
   created_time       2023-05-10T14:32:32.37354939Z
   custom_metadata    <nil>
   deletion_time      n/a
   destroyed          false
   version            1

   ====== Data ======
   Key         Value
   ---         -----
   password    !hh&be1e4j16dVc0ggae
   ```

### Deleting and destroying data

When deleting data the standard `d8 stronghold kv delete` command will perform a
soft delete. It will mark the version as deleted and populate a `deletion_time`
timestamp. Soft deletes do not remove the underlying version data from storage,
which allows the version to be undeleted. The `d8 stronghold kv undelete` command
handles undeleting versions.

A version's data is permanently deleted only when the key has more versions than
are allowed by the max-versions setting, or when using `d8 stronghold kv destroy`. When
the destroy command is used the underlying version data will be removed and the
key metadata will be marked as destroyed. If a version is cleaned up by going
over max-versions the version metadata will also be removed from the key.

See the commands below for more information:

1. The latest version of a key can be deleted with the delete command, this also
   takes a `-versions` flag to delete prior versions:

   ```shell-session
   $ d8 stronghold kv delete -mount=secret my-secret
   Success! Data deleted (if it existed) at: secret/data/my-secret
   ```

1. Versions can be undeleted:

   ```shell-session
   $ d8 stronghold kv undelete -mount=secret -versions=2 my-secret
   Success! Data written to: secret/undelete/my-secret

   $ d8 stronghold kv get -mount=secret my-secret
   ====== Metadata ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:23:21.834403Z
   custom_metadata  <nil>
   deletion_time    n/a
   destroyed        false
   version          2

   ====== Data ======
   Key         Value
   ---         -----
   my-value    short-lived-s3cr3t
   ```

1. Destroying a version permanently deletes the underlying data:

   ```shell-session
   $ d8 stronghold kv destroy -mount=secret -versions=2 my-secret
   Success! Data written to: secret/destroy/my-secret
   ```

### Key metadata

All versions and key metadata can be tracked with the metadata command & API.
Deleting the metadata key will cause all metadata and versions for that key to
be permanently removed.

See the commands below for more information:

1. All metadata and versions for a key can be viewed:

   ```shell-session
   $ d8 stronghold kv metadata get -mount=secret my-secret
   ========== Metadata ==========
   Key                     Value
   ---                     -----
   cas_required            false
   created_time            2024-06-19T17:20:22.985303Z
   current_version         2
   custom_metadata         <nil>
   delete_version_after    0s
   max_versions            0
   oldest_version          0
   updated_time            2024-06-19T17:22:23.369372Z

   ====== Version 1 ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:20:22.985303Z
   deletion_time    n/a
   destroyed        false

   ====== Version 2 ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:22:23.369372Z
   deletion_time    n/a
   destroyed        true
   ```

1. The metadata settings for a key can be configured:

   ```shell-session
   $ d8 stronghold kv metadata put -mount=secret -max-versions 2 -delete-version-after="3h25m19s" my-secret
   Success! Data written to: secret/metadata/my-secret
   ```

   Delete-version-after settings will apply only to new versions. Max versions
   changes will be applied on next write:

   ```shell-session
   $ d8 stronghold kv put -mount=secret my-secret my-value=newer-s3cr3t
   Key              Value
   ---              -----
   created_time     2024-06-19T17:31:16.662563Z
   custom_metadata  <nil>
   deletion_time    2024-06-19T20:56:35.662563Z
   destroyed        false
   version          4
   ```

   Once a key has more versions than max versions the oldest versions
   are cleaned up:

   ```shell-session
   $ d8 stronghold kv metadata get -mount=secret my-secret
   ========== Metadata ==========
   Key                     Value
   ---                     -----
   cas_required            false
   created_time            2024-06-19T17:20:22.985303Z
   current_version         4
   custom_metadata         <nil>
   delete_version_after    3h25m19s
   max_versions            2
   oldest_version          3
   updated_time            2024-06-19T17:31:16.662563Z

   ====== Version 3 ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:23:21.834403Z
   deletion_time    n/a
   destroyed        true

   ====== Version 4 ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:31:16.662563Z
   deletion_time    2024-06-19T20:56:35.662563Z
   destroyed        false
   ```

   A secret's key metadata can contain custom metadata used to describe the secret.
   The data will be stored as string-to-string key-value pairs.
   The `-custom-metadata` flag can be repeated to add multiple key-value pairs.

   The `d8 stronghold kv metadata put` command can be used to fully overwrite the value of `custom_metadata`:

   ```shell-session
   $ d8 stronghold kv metadata put -mount=secret -custom-metadata=foo=abc -custom-metadata=bar=123 my-secret
   Success! Data written to: secret/metadata/my-secret

   $ d8 stronghold kv get -mount=secret my-secret
   ====== Metadata ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:22:23.369372Z
   custom_metadata  map[bar:123 foo:abc]
   deletion_time    n/a
   destroyed        false
   version          2

   ====== Data ======
   Key         Value
   ---         -----
   foo         aa
   bar         bb
   ```

   The `d8 stronghold kv metadata patch` command can be used to partially overwrite the value of `custom_metadata`.
   The following invocation will update `custom_metadata` sub-field `foo` but leave `bar` untouched:

   ```shell-session
   $ d8 stronghold kv metadata patch -mount=secret -custom-metadata=foo=def my-secret
   Success! Data written to: secret/metadata/my-secret
   ```

   ```shell-session
   $ d8 stronghold kv get -mount=secret my-secret
   ====== Metadata ======
   Key              Value
   ---              -----
   created_time     2024-06-19T17:22:23.369372Z
   custom_metadata  map[bar:123 foo:def]
   deletion_time    n/a
   destroyed        false
   version          2

   ====== Data ======
   Key         Value
   ---         -----
   foo         aa
   bar         bb
   ```

1. Permanently delete all metadata and versions for a key:

   ```shell-session
   $ d8 stronghold kv metadata delete -mount=secret my-secret
   Success! Data deleted (if it existed) at: secret/metadata/my-secret
   ```

## API

The KV secrets engine has a full HTTP API. Please see the
[KV secrets engine API](/api-docs/secret/kv/kv-v2) for more
details.
