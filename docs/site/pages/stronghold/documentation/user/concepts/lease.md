---
title: "Lease"
permalink: en/stronghold/documentation/user/concepts/lease.html
lang: en
description: >-
  Stronghold provides a lease with every secret. When this lease is expired, Stronghold
  will revoke that secret.
---

## Lease, renew, and revoke

With every dynamic secret and `service` type authentication token, Stronghold
creates a _lease_: metadata containing information such as a time duration,
renewability, and more. Stronghold promises that the data will be valid for the
given duration, or Time To Live (TTL). Once the lease is expired, Stronghold can
automatically revoke the data, and the consumer of the secret can no longer be
certain that it is valid.

The benefit should be clear: consumers of secrets need to check in with
Stronghold routinely to either renew the lease (if allowed) or request a
replacement secret. This makes the Stronghold audit logs more valuable and
also makes key rolling a lot easier.

All dynamic secrets in Stronghold are required to have a lease. Even if the data is
meant to be valid for eternity, a lease is required to force the consumer
to check in routinely.

In addition to renewals, a lease can be _revoked_. When a lease is revoked, it
invalidates that secret immediately and prevents any further renewals. For
example, with the [Kubernetes secrets engine](/docs/secrets/kubernetes), the
service accounts will be deleted from Kubernetes the moment a lease is revoked. This
renders the access keys invalid from that point forward.

Revocation can happen manually via the API, via the `d8 stronghold lease revoke` cli command,
the user interface (UI) under the Access tab, or automatically by Stronghold. When a lease
is expired, Stronghold will automatically revoke that lease. When a token is revoked,
Stronghold will revoke all leases that were created using that token.

**Note**: The [Key/Value Backend](/docs/secrets/kv) which stores
arbitrary secrets does not issue leases although it will sometimes return a
lease duration; see the documentation for more information.

## Lease IDs

When reading a dynamic secret, such as via `d8 stronghold read`, Stronghold always returns a
`lease_id`. This is the ID used with commands such as `d8 stronghold lease renew` and `d8 stronghold lease revoke` to manage the lease of the secret.

## Lease durations and renewal

Along with the lease ID, a _lease duration_ can be read. The lease duration is
a Time To Live value: the time in seconds for which the lease is valid. A
consumer of this secret must renew the lease within that time.

When renewing the lease, the user can request a specific amount of time they
want remaining on the lease, termed the `increment`. This is not an increment
at the end of the current TTL; it is an increment _from the current time_. For
example, `d8 stronghold lease renew -increment=3600 my-lease-id` would request that the TTL of the lease
be adjusted to 1 hour (3600 seconds). Having the increment be rooted at the
current time instead of the end of the lease makes it easy for users to reduce
the length of leases if they don't actually need credentials for the full
possible lease period, allowing those credentials to expire sooner and
resources to be cleaned up earlier.

The requested increment is completely advisory. The backend in charge of the
secret can choose to completely ignore it. For most secrets, the backend does
its best to respect the increment, but often limits it to ensure renewals every
so often.

As a result, the return value of renewals should be carefully inspected to
determine what the new lease is.

{% alert level="info" %}

To implement token renewal logic in your application code, refer to the [code example in the Authentication doc](/docs/concepts/auth#code-example).

{% endalert %}

## Prefix-based revocation

In addition to revoking a single secret, operators with proper access control
can revoke multiple secrets based on their lease ID prefix.

Lease IDs are structured in a way that their prefix is always the path where
the secret was requested from. This lets you revoke trees of secrets. For
example, to revoke all Userpass logins, you can do `d8 stronghold lease revoke -prefix auth/userpass/`.
For more information about revoke command please check
[cli's lease revoke](/docs/commands/lease/revoke#lease-revoke)
command docs.

This is very useful if there is an intrusion within a specific system: all
secrets of a specific backend or a certain configured backend can be revoked
quickly and easily.
