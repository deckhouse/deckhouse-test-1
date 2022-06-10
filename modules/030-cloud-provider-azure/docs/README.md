---
title: "Cloud provider — Azure"
---

The `cloud-provider-azure` module:
- Manages Azure resources using the `cloud-controller-manager` (CCM) module:
    * The CCM module creates network routes for the `PodNetwork` network on the Azure side;
    * The CCM module creates LoadBalancers for Kubernetes Service objects that have the `LoadBalancer` type;
    * The CCM module updates the metadata of the cluster nodes according to the configuration parameters and deletes nodes that are no longer in Azure;
- Provisions nodes in Azure using the `CSI storage` component;
- Enables the necessary CNI plugin (using the [simple bridge](../../modules/035-cni-simple-bridge/));
- Registers with the [node-manager](../../modules/040-node-manager/) module so that [AzureInstanceClasses](cr.html#azureinstanceclass) can be used when creating the [NodeGroup](../../modules/040-node-manager/cr.html#nodegroup).
