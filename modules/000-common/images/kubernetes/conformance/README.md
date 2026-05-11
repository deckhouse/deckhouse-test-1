# Reproducing the test results

1) Deploy Deckhouse
2) Disable the security restrictions
```bash
kubectl apply -f - <<"EOF"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:kubelet-api-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:kubelet-api-admin
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: kube-apiserver-kubelet-client
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleConfig
metadata:
  name: admission-policy-engine
spec:
  enabled: true
  settings:
    podSecurityStandards:
      defaultPolicy: Privileged
  version: 1
EOF
```

3) Download a [binary release](https://github.com/vmware-tanzu/sonobuoy/releases) of the CLI.
```bash
wget https://github.com/vmware-tanzu/sonobuoy/releases/download/v0.57.3/sonobuoy_0.57.3_linux_amd64.tar.gz && tar xzf sonobuoy_0.57.3_linux_amd64.tar.gz
```

4) Run the tests
```bash
./sonobuoy run --mode=certified-conformance
```

5) Once `./sonobuoy status` shows the run as `completed`, copy the output directory from the main Sonobuoy pod to 
```bash
./sonobuoy retrieve . -f sb.tar.gz && tar -xzf sb.tar.gz --strip-components=4 -C . plugins/e2e/results/global/e2e.log plugins/e2e/results/global/junit_01.xml && mv plugins/e2e/results/global/e2e.log . && mv plugins/e2e/results/global/junit_01.xml .
```

6) Copy these 2 files in a version directory
7) Place the `tests/conformance/$version` label to the PR
