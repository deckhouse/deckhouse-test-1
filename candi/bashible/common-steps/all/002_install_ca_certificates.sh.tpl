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

bb-package-install "d8-ca-updater:{{ .images.registrypackages.d8CaUpdater001 }}"

# hack to avoid problems with certs in alpine busybox for kube-apiserver
if [[ ! -e /etc/ssl/certs/ca-certificates.crt ]]; then
  ln -s /etc/ssl/certs/ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
fi

{{- if .registry.ca }}

bb-sync-file /usr/local/share/d8-ca-certificates/mozilla/registry-ca.crt - registry-ca-changed << "EOF"
{{ .registry.ca }}
EOF

bb-flag-set containerd-need-restart

{{- else }}
if [ -f /usr/local/share/d8-ca-certificates/mozilla/registry-ca.crt ]; then
  rm -f /usr/local/share/d8-ca-certificates/mozilla/registry-ca.crt
fi
{{- end }}

/opt/deckhouse/bin/d8-ca-updater
