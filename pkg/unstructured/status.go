package unstructured

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//GetCondition returns the condition with type typeString
// returns error if the condition is not found
func GetCondition(u *unstructured.Unstructured, typeString string) (map[string]interface{}, error) {
	if u != nil {
		if v, ok := u.Object["status"]; ok {
			status := v.(map[string]interface{})
			if v, ok := status["conditions"]; ok {
				conditions := v.([]interface{})
				for _, v := range conditions {
					condition := v.(map[string]interface{})
					if v, ok := condition["type"]; ok {
						if v.(string) == typeString {
							return condition, nil
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("condition %s not found", typeString)
}
