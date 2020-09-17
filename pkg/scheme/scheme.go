package scheme

import (
	"sync"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/scale/scheme"
	"k8s.io/klog"
)

var k8sNativeScheme *runtime.Scheme
var k8sNativeSchemeOnce sync.Once

func KubernetesNativeScheme() *runtime.Scheme {
	k8sNativeSchemeOnce.Do(func() {
		k8sNativeScheme = runtime.NewScheme()
		err := scheme.AddToScheme(k8sNativeScheme)
		if err != nil {
			klog.Error(err)
			return
		}
		// API extensions are not in the above scheme set,
		// and must thus be added separately.
		err = apiextensionsv1beta1.AddToScheme(k8sNativeScheme)
		if err != nil {
			klog.Error(err)
			return
		}
		err = apiextensionsv1.AddToScheme(k8sNativeScheme)
		if err != nil {
			klog.Error(err)
			return
		}
	})
	return k8sNativeScheme
}
