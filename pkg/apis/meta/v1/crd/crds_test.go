package crd

import (
	"testing"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	fakeclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHaveCRDs(t *testing.T) {

	crd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "Test",
		},
	}

	client := fakeclientset.NewSimpleClientset(crd)

	type args struct {
		client       clientset.Interface
		expectedCRDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "check if crd present",
			args: args{
				client:       client,
				expectedCRDs: []string{"Test"},
			},
			wantErr: false,
		},
		{
			name: "check all crds not present",
			args: args{
				client:       client,
				expectedCRDs: []string{"Test", "Test2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HaveCRDs(tt.args.client, tt.args.expectedCRDs); (err != nil) != tt.wantErr {
				t.Errorf("HaveCRDs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
