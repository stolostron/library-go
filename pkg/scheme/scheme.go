package scheme

import (
	"sync"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/scale/scheme"
	"k8s.io/klog"
)

var k8sNativeScheme *runtime.Scheme
var k8sNativeSchemeOnce sync.Once

// AsVersioned converts the given info into a runtime.Object with the correct
// group and version set
// func AsVersioned(info *resource.Info) runtime.Object {
// 	return convertWithMapper(info.Object, info.Mapping)
// }

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
		err = corev1.AddToScheme(k8sNativeScheme)
		if err != nil {
			klog.Error(err)
			return
		}
	})
	return k8sNativeScheme
}

// ConvertWithMapper converts the given object with the optional provided
// RESTMapping. If no mapping is provided, the default schema versioner is used
func ConvertWithMapper(obj runtime.Object, mapping *meta.RESTMapping) runtime.Object {
	s := KubernetesNativeScheme()
	var gv = runtime.GroupVersioner(schema.GroupVersions(s.PrioritizedVersionsAllGroups()))
	if mapping != nil {
		gv = mapping.GroupVersionKind.GroupVersion()
	}
	if obj, err := runtime.ObjectConvertor(s).ConvertToVersion(obj, gv); err == nil {
		return obj
	}
	return obj
}
