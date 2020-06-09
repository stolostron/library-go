// Copyright (c) 2020 Red Hat, Inc.

package applier

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestApplierClient_CreateOrUpdateInPath(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	a, err := NewApplier(NewTestReader(), values, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	client := fake.NewFakeClient([]runtime.Object{}...)

	c, err := NewApplierClient(a, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
	}
	tests := []struct {
		name    string
		fields  Client
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: *c,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewApplierClient(tt.fields.applier, tt.fields.client, tt.fields.owner, tt.fields.scheme, tt.fields.merger)
			if err != nil {
				t.Error(err)
			}
			if err := c.CreateOrUpdateInPath(tt.args.path, tt.args.excluded, tt.args.recursive); (err != nil) != tt.wantErr {
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
