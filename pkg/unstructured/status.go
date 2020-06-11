package unstructured

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

//StatusContainsTypeEqualTo check if u contains a condition type with value typeString
func StatusContainsTypeEqualTo(u *unstructured.Unstructured, typeString string) bool {
	if u != nil {
		if v, ok := u.Object["status"]; ok {
			status := v.(map[string]interface{})
			if v, ok := status["conditions"]; ok {
				conditions := v.([]interface{})
				for _, v := range conditions {
					condition := v.(map[string]interface{})
					if v, ok := condition["type"]; ok {
						if v.(string) == typeString {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
