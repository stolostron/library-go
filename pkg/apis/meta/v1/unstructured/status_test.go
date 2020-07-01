package unstructured

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestGetConditionByType(t *testing.T) {
	u := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cluster.open-cluster-management.io/v1",
			"kind":       "ManagedCluster",
			"metadata": map[string]interface{}{
				"name": "myname",
			},
			"status": map[string]interface{}{
				"conditions": []interface{}{
					map[string]interface{}{
						"type":   "mytype",
						"reason": "test",
					},
				},
			},
		},
	}
	unostatus := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cluster.open-cluster-management.io/v1",
			"kind":       "ManagedCluster",
			"metadata": map[string]interface{}{
				"name": "myname",
			},
		},
	}
	type args struct {
		u          *unstructured.Unstructured
		typeString string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "condition found",
			args: args{
				u:          u,
				typeString: "mytype",
			},
			want: map[string]interface{}{
				"type":   "mytype",
				"reason": "test",
			},
			wantErr: false,
		},
		{
			name: "condition not found",
			args: args{
				u:          u,
				typeString: "notExists",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "status not found",
			args: args{
				u:          unostatus,
				typeString: "notExists",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unstructured is nil",
			args: args{
				u:          nil,
				typeString: "notExists",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConditionByType(tt.args.u, tt.args.typeString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
