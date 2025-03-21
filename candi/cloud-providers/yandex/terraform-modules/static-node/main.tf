# Copyright 2021 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

locals {
  mapping = lookup(var.providerClusterConfiguration, "existingZoneToSubnetIDMap", {})

  zone_to_subnet = length(local.mapping) == 0 ? {
    "ru-central1-a" = length(data.yandex_vpc_subnet.kube_a) > 0 ? data.yandex_vpc_subnet.kube_a[0] : object({})
    "ru-central1-b" = length(data.yandex_vpc_subnet.kube_b) > 0 ? data.yandex_vpc_subnet.kube_b[0] : object({})
    "ru-central1-d" = length(data.yandex_vpc_subnet.kube_d) > 0 ? data.yandex_vpc_subnet.kube_d[0] : object({})
  } : data.yandex_vpc_subnet.existing

  actual_zones    = lookup(var.providerClusterConfiguration, "zones", null) != null ? tolist(setintersection(keys(local.zone_to_subnet), var.providerClusterConfiguration.zones)) : keys(local.zone_to_subnet)
  zones           = lookup(local.ng, "zones", null) != null ? tolist(setintersection(local.actual_zones, local.ng["zones"])) : local.actual_zones
  subnets         = length(local.zones) > 0 ? [for z in local.zones : local.zone_to_subnet[z]] : values(local.zone_to_subnet)
  internal_subnet = element(local.subnets, var.nodeIndex)

  // TODO apply external_subnet_id_from_ids to external_subnet_id directly after remove externalSubnetID
  external_subnet_id_from_ids = length(local.external_subnet_ids) > 0 ? element(local.external_subnet_ids, var.nodeIndex) : null

  external_subnet_id         = local.external_subnet_id_from_ids == null ? local.external_subnet_id_deprecated : local.external_subnet_id_from_ids
  assign_external_ip_address = (local.external_subnet_id == null) && (length(local.external_ip_addresses) > 0) ? true : false
}

data "yandex_vpc_subnet" "existing" {
  for_each = local.mapping
  subnet_id = each.value
}

data "yandex_vpc_subnet" "kube_a" {
  count = length(local.mapping) == 0 ? 1 : 0
  name = "${local.prefix}-a"
}

data "yandex_vpc_subnet" "kube_b" {
  count = length(local.mapping) == 0 ? 1 : 0
  name = "${local.prefix}-b"
}

data "yandex_vpc_subnet" "kube_d" {
  count = length(local.mapping) == 0 ? 1 : 0
  name = "${local.prefix}-d"
}

resource "yandex_vpc_address" "addr" {
  count = var.nodeIndex < length(local.external_ip_addresses) ? local.external_ip_addresses[var.nodeIndex] == "Auto" ? 1 : 0 : 0
  name  = join("-", [local.prefix, var.nodeGroupName, var.nodeIndex])

  external_ipv4_address {
    zone_id = local.internal_subnet.zone
  }
}

locals {
  # null if var.nodeIndex < length(local.external_ip_addresses)
  # yandex_vpc_address.addr[0].external_ipv4_address[0].address if local.external_ip_addresses == Auto
  # local.external_ip_addresses[var.nodeIndex] if local.external_ip_addresses contain IP-addresses
  external_ip = var.nodeIndex < length(local.external_ip_addresses) ? local.external_ip_addresses[var.nodeIndex] == "Auto" ? yandex_vpc_address.addr[0].external_ipv4_address[0].address : local.external_ip_addresses[var.nodeIndex] : null
}

resource "yandex_compute_instance" "static" {
  name     = join("-", [local.prefix, var.nodeGroupName, var.nodeIndex])
  hostname = join("-", [local.prefix, var.nodeGroupName, var.nodeIndex])
  zone     = local.internal_subnet.zone

  allow_stopping_for_update = true

  platform_id = local.platform

  resources {
    cores         = local.cores
    core_fraction = local.core_fraction
    memory        = local.memory
  }

  labels = local.additional_labels

  boot_disk {
    initialize_params {
      type     = local.disk_type
      image_id = local.image_id
      size     = local.disk_size_gb
    }
  }

  dynamic "network_interface" {
    for_each = local.external_subnet_id != null ? [local.external_subnet_id] : []
    content {
      subnet_id = network_interface.value
      nat       = false
    }
  }

  network_interface {
    subnet_id      = local.internal_subnet.id
    nat            = local.assign_external_ip_address
    nat_ip_address = local.assign_external_ip_address ? local.external_ip : null
  }

  network_acceleration_type = local.network_type

  lifecycle {
    ignore_changes = [
      metadata,
      secondary_disk,
    ]
  }

  timeouts {
    create = var.resourceManagementTimeout
    delete = var.resourceManagementTimeout
    update = var.resourceManagementTimeout
  }

  metadata = {
    ssh-keys          = "user:${local.ssh_public_key}"
    user-data         = base64decode(var.cloudConfig)
    node-network-cidr = local.node_network_cidr
  }
}
