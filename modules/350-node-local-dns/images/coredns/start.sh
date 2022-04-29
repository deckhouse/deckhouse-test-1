#!/bin/bash

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

set -Eeuo pipefail

# Setup interface
dev_name="nodelocaldns"
if ! ip link show "$dev_name" >/dev/null 2>&1 ; then
  ip link add "$dev_name" type dummy
fi
if ! ip -json addr show "$dev_name" | jq -re "any(.[].addr_info[]?.local; . == \"169.254.20.10\")" >/dev/null 2>&1 ; then
  ip addr add 169.254.20.10/32 dev "$dev_name"
fi
if ! ip -json addr show "$dev_name" | jq -re "any(.[].addr_info[]?.local; . == \"${KUBE_DNS_SVC_IP}\")" >/dev/null 2>&1 ; then
  ip addr add "${KUBE_DNS_SVC_IP}"/32 dev "$dev_name"
fi

# Setup iptables
if ! iptables -w 60 -t raw -C OUTPUT -s "${KUBE_DNS_SVC_IP}/32" -p tcp -m tcp --sport 53 -j NOTRACK >/dev/null 2>&1; then
  iptables -w 60 -t raw -A OUTPUT -s "${KUBE_DNS_SVC_IP}/32" -p tcp -m tcp --sport 53 -j NOTRACK
fi
if ! iptables -w 60 -t raw -C OUTPUT -s "${KUBE_DNS_SVC_IP}/32" -p udp -m udp --sport 53 -j NOTRACK >/dev/null 2>&1; then
  iptables -w 60 -t raw -A OUTPUT -s "${KUBE_DNS_SVC_IP}/32" -p udp -m udp --sport 53 -j NOTRACK
fi
if iptables -w 60 -t raw -C PREROUTING -d "${KUBE_DNS_SVC_IP}/32" -m socket --nowildcard -j NOTRACK >/dev/null 2>&1 ; then
  # Remove. Will be added later.
  iptables -w 60 -t raw -D PREROUTING -d "${KUBE_DNS_SVC_IP}/32" -m socket --nowildcard -j NOTRACK >/dev/null 2>&1
fi

echo -n "not-ready" > /tmp/coredns-readiness

exec /coredns -conf /etc/coredns/Corefile
