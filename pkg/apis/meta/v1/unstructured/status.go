package unstructured

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//GetConditionByType returns the condition with type typeString
// returns error if the condition is not found
//u: The *unstructured.Unstructured object to search in
//typeString: the type to search
func GetConditionByType(
	u *unstructured.Unstructured,
	typeString string,
) (map[string]interface{}, error) {
	if u != nil {
		if v, ok := u.Object["status"]; ok {
			status := v.(map[string]interface{})
			if v, ok := status["conditions"]; ok {
				return searchCondition(v.([]interface{}), typeString)
			}
		}
		return nil, fmt.Errorf("status not found")
	}
	return nil, fmt.Errorf("the passed unstructured is nil")
}

func searchCondition(
	conditions []interface{},
	typeString string,
) (map[string]interface{}, error) {
	for _, v := range conditions {
		condition := v.(map[string]interface{})
		if v, ok := condition["type"]; ok {
			if v.(string) == typeString {
				return condition, nil
			}
		}
	}
	return nil, fmt.Errorf("condition %s not found", typeString)
}
