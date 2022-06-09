---
title: "The deckhouse module"
search: releaseChannel, release channel stabilization, auto-switching the release channel
---

Every release can have one of the following statuses:
  * `Pending` - release is waiting to be deployed: waiting for update window, canary deployment, etc. You can see the detailed status via the `kubectl describe deckhouserelease $name` command.
  * `Deployed` - release is applied. It means that the Deckhouse image tag was changed, but the update process of all components
is going asynchronously and could not have been finished yet.
  * `Outdated` - release is outdated and not used anymore.
  * `Suspended` - release was suspended (for ex. it has an error). Can be set only if `suspended` release was not deployed yet.


#### Update process
When release status is changed to `Deployed` state, release is updating only a tag of the Deckhouse image.
Deckhouse will start checking and updating process of the all modules, which were changed from the last release.
Duration of an update could be different and connected to cluster size, enabled modules count and settings.
Foe example: if you cluster have a lot of `NodeGroup` resources, it will take some time to update them because these resources are updated one by one
`IngressNginxControllers` also updating one by one.


#### Manual release deployment
If you have a [manual update mode](usage.html#manual-update-confirmation) enabled and have a few Pending releases,
you can approve them all at once. In that case Deckhouse will update in series keeping a release order and changing their status during the update.


#### Pin a release
Release pinning could be necessary if you want to hold a Deckhouse update for some reason.

There are 3 options to pin a release:
- Set a [manual update mode](usage.html#manual-update-confirmation).
In this case, you will hold a current release but patch-release will still be applied. Minor-release will not be changed without your approval.

  Example:
    The current release is `v1.29.3`, after setting a manual update mode Deckhouse will be able to update to version `v1.29.9` but won't be able to apply version `v1.30.0`.

- Set a specified image tag for deployment/deckhouse. 
In this case, you will hold a Deckhouse version until a new release will come.
You may need this in a situation when some Deckhouse release has an error that hasn't occurred earlier and you want to roll back to the previous release but update as soon as a new release with a patch will come.

  Example:
    `kubectl -n d8-system set image deployment/deckhouse deckhouse=registry.deckhouse.io/deckhouse/fe:v1.30.5`

- Set a specified image tag for deployment/deckhouse and remove `releaseChannel` from deckhouse ConfigMap.
    In this case, you will hold a specified version and will not get any more updates.
    ```sh
    $ kubectl -n d8-system set image deployment/deckhouse deckhouse=registry.deckhouse.io/deckhouse/fe:v1.30.5
    $ kubectl -n d8-system edit cm deckhouse
      // remove releaseChannel
    ```
