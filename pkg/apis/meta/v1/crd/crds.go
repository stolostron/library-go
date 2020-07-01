package crd

import (
	"context"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

//HaveCRDs returns an error if all provided CRDs are not installed
//client: the client to use
//expectedCRDs: The list of expected CRDS to find
func HaveCRDs(client clientset.Interface, expectedCRDs []string) error {
	clientAPIExtensionV1beta1 := client.ApiextensionsV1beta1()
	for _, crd := range expectedCRDs {
		klog.V(1).Infof("Check if %s exists", crd)
		_, err := clientAPIExtensionV1beta1.CustomResourceDefinitions().Get(context.TODO(), crd, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving crd %s: %s", crd, err.Error())
			return err
		}
	}
	return nil
}
