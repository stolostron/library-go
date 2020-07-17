module github.com/open-cluster-management/library-go

go 1.13

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/tools v0.0.0-20200713235242-6acd2ab80ede // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	k8s.io/api v0.18.4
	k8s.io/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cluster-registry v0.0.6
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.18.4 // indirect
	sigs.k8s.io/controller-runtime v0.6.0
)

replace k8s.io/client-go => k8s.io/client-go v0.18.2
