// Copyright Contributors to the Open Cluster Management project

package crd

import (
	"reflect"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	fakeclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHasCRDs(t *testing.T) {

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
		want    bool
		want1   []string
		wantErr bool
	}{
		{
			name: "check if crd present",
			args: args{
				client:       client,
				expectedCRDs: []string{"Test"},
			},
			want:    true,
			want1:   []string{},
			wantErr: false,
		},
		{
			name: "check all crds not present",
			args: args{
				client:       client,
				expectedCRDs: []string{"Test", "Test2"},
			},
			want:    false,
			want1:   []string{"Test2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := HasCRDs(tt.args.client, tt.args.expectedCRDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasCRDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HasCRDs() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("HasCRDs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
