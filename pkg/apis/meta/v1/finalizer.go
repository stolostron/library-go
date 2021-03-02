// Copyright Contributors to the Open Cluster Management project

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//AddFinalizer accepts cluster and adds provided finalizer to cluster
func AddFinalizer(o metav1.Object, finalizer string) {
	for _, f := range o.GetFinalizers() {
		if f == finalizer {
			return
		}
	}

	o.SetFinalizers(append(o.GetFinalizers(), finalizer))
}

//RemoveFinalizer accepts cluster and removes provided finalizer if present
func RemoveFinalizer(o metav1.Object, finalizer string) {
	var finalizers []string

	for _, f := range o.GetFinalizers() {
		if f != finalizer {
			finalizers = append(finalizers, f)
		}
	}

	if len(finalizers) == len(o.GetFinalizers()) {
		return
	}

	o.SetFinalizers(finalizers)
}
