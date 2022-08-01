---
title: "The monitoring-ping module"
---

## Description

This module monitors network connectivity between cluster nodes and external nodes (optionally).

## How does it work?

The module tracks the node's `.status.addresses` field for changes. Upon detecting changes, it invokes a hook that collects a complete list of node names/addresses and passes it to a DaemonSet (the latter recreates the Pods). As a result, ping checks the always up-to-date list of nodes.
