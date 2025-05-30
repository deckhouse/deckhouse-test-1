# Copyright 2021 Flant CJSC
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

module "static-node" {
  source = "../../../terraform-modules/static-node"
  prefix = local.prefix
  cluster_uuid = var.clusterUUID
  node_index = var.nodeIndex
  node_group = local.node_group
  root_volume_size = local.root_volume_size
  root_volume_type = local.root_volume_type
  additional_security_groups = local.additional_security_groups
  cloud_config = var.cloudConfig
  zones = local.zones
  tags = local.tags
  resourceManagementTimeout = var.resourceManagementTimeout
}
