module github.com/open-cluster-management/library-go

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	k8s.io/api v0.18.3
	k8s.io/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.0
)

replace k8s.io/client-go => k8s.io/client-go v0.18.2
