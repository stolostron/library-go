module github.com/open-cluster-management/library-go

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	github.com/sclevine/agouti v3.0.0+incompatible
	k8s.io/api v0.18.3
	k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.4.0
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
)
