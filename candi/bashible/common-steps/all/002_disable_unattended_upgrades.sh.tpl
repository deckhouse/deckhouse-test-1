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

if ! systemctl list-unit-files unattended-upgrades.service >/dev/null 2>&1; then
  exit 0
fi

if systemctl is-enabled --quiet unattended-upgrades ; then
  systemctl disable --now unattended-upgrades
fi

bb-apt-remove unattended-upgrades
