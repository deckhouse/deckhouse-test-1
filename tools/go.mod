module tools

go 1.15

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/golangci/golangci-lint v1.40.1
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	golang.org/x/tools v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/klog/v2 v2.100.1 // indirect
	sigs.k8s.io/controller-tools v0.0.0
	sigs.k8s.io/yaml v1.3.0
// k8s.io/apiextensions-apiserver v0.0.0
)

replace (
	k8s.io/apiextensions-apiserver v0.0.0 => github.com/alex123012/deckhouse-controller-tools/pkg/apiextensions-apiserver v0.0.0-20230510090815-d594daf1af8c
	sigs.k8s.io/controller-tools v0.0.0 => github.com/alex123012/deckhouse-controller-tools v0.0.0-20230510090815-d594daf1af8c
)
