package applier

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestApplierClient_CreateOrUpdateInPath(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	client := fake.NewFakeClient([]runtime.Object{}...)

	a, err := NewApplier(tp, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
		values    interface{}
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: *a,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    values,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.CreateOrUpdateInPath(tt.args.path, tt.args.excluded, tt.args.recursive, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("ApplierClient.CreateOrUpdateInPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				sa := &corev1.ServiceAccount{}
				err := client.Get(context.TODO(), types.NamespacedName{
					Name:      values.ManagedClusterName,
					Namespace: values.ManagedClusterNamespace,
				}, sa)
				if err != nil {
					t.Error(err)
				}
				r := &rbacv1.ClusterRole{}
				err = client.Get(context.TODO(), types.NamespacedName{
					Name: values.ManagedClusterName,
				}, r)
				if err != nil {
					t.Error(err)
				}
				rb := &rbacv1.ClusterRoleBinding{}
				err = client.Get(context.TODO(), types.NamespacedName{
					Name: values.ManagedClusterName,
				}, rb)
				if err != nil {
					t.Error(err)
				}
				if rb.RoleRef.Name != "system:test:"+values.ManagedClusterName {
					t.Errorf("Expecting %s got %s", "system:test:"+values.ManagedClusterName, rb.RoleRef.Name)
				}
			}
		})
	}
}

func TestNewApplier(t *testing.T) {
	tp := &TemplateProcessor{}
	client := fake.NewFakeClient([]runtime.Object{}...)
	owner := &corev1.Secret{}
	scheme := &runtime.Scheme{}
	merger := func(current,
		new *unstructured.Unstructured,
	) (
		future *unstructured.Unstructured,
		update bool,
	) {
		return nil, true
	}
	type args struct {
		templateProcessor *TemplateProcessor
		client            crclient.Client
		owner             metav1.Object
		scheme            *runtime.Scheme
		merger            Merger
	}
	tests := []struct {
		name    string
		args    args
		want    *Applier
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				templateProcessor: tp,
				client:            client,
				owner:             owner,
				scheme:            scheme,
				merger:            merger,
			},
			want: &Applier{
				templateProcessor: tp,
				client:            client,
				owner:             owner,
				scheme:            scheme,
				merger:            merger,
			},
			wantErr: false,
		},
		{
			name: "Failed no templateProcessor",
			args: args{
				templateProcessor: nil,
				client:            client,
				owner:             owner,
				scheme:            scheme,
				merger:            merger,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Failed no client",
			args: args{
				templateProcessor: tp,
				client:            nil,
				owner:             owner,
				scheme:            scheme,
				merger:            merger,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApplier(tt.args.templateProcessor, tt.args.client, tt.args.owner, tt.args.scheme, tt.args.merger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.templateProcessor, tt.want.templateProcessor) &&
					!reflect.DeepEqual(got.client, tt.want.client) &&
					!reflect.DeepEqual(got.owner, tt.want.owner) &&
					!reflect.DeepEqual(got.scheme, tt.want.scheme) &&
					!reflect.DeepEqual(got.merger, tt.want.merger) {
					t.Errorf("NewApplier() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
