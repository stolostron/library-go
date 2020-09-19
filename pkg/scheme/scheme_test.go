package scheme

import (
	"reflect"
	"testing"

	"gopkg.in/square/go-jose.v2/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

var values = struct {
	ManagedClusterName          string
	ManagedClusterNamespace     string
	BootstrapServiceAccountName string
}{
	ManagedClusterName:          "mycluster",
	ManagedClusterNamespace:     "myclusterns",
	BootstrapServiceAccountName: "mysa",
}

func TestConvertWithMapper(t *testing.T) {
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      values.BootstrapServiceAccountName,
			Namespace: values.ManagedClusterNamespace,
		},
	}

	jsa, err := json.Marshal(sa)
	if err != nil {
		t.Error(err)
	}
	usa := &unstructured.Unstructured{}
	err = usa.UnmarshalJSON(jsa)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		obj     runtime.Object
		mapping *meta.RESTMapping
	}
	tests := []struct {
		name string
		args args
		want runtime.Object
	}{
		{
			name: "Succeed",
			args: args{
				obj:     usa,
				mapping: nil,
			},
			want: sa,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertWithMapper(tt.args.obj, tt.args.mapping); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertWithMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}
