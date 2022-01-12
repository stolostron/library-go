// Copyright Contributors to the Open Cluster Management project

package unstructured

import (
	libgounstructuredv1 "github.com/stolostron/library-go/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//Deprecated:
// Use https://github.com/stolostron/library-go/pkg/apis/meta/v1/unstructured#GetCondition
//GetCondition returns the condition with type typeString
// returns error if the condition is not found
//u: The *unstructured.Unstructured object to search in
//typeString: the type to search
func GetCondition(
	u *unstructured.Unstructured,
	typeString string,
) (map[string]interface{}, error) {
	return libgounstructuredv1.GetConditionByType(u, typeString)
}
