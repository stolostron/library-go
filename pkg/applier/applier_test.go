// Copyright (c) 2020 Red Hat, Inc.

package applier

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestTemplateAssetsToMapOfUnstructured(t *testing.T) {
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	type config struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}
	type args struct {
		path      string
		recursive bool
		config    config
	}
	type check struct {
		kinds []string
	}
	tests := []struct {
		name       string
		args       args
		check      check
		wantAssets map[string]*unstructured.Unstructured
		wantErr    bool
	}{
		{
			name: "Parse",
			args: args{
				path:      "test",
				recursive: true,
				config: config{
					ManagedClusterName:          "mymanagedcluster",
					ManagedClusterNamespace:     "mymanagedclusterNS",
					BootstrapServiceAccountName: "mymanagedcluster",
				},
			},
			check: check{
				kinds: []string{
					"ClusterRole",
					"ClusterRoleBinding",
					"ServiceAccount",
					"Namespace",
					"Secret",
					"ServiceAccount",
					"Role",
					"RoleBinding",
					"ConfigMap",
					"Deployment",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAssets, err := a.TemplateAssetsInPathUnstructured(tt.args.path, nil, tt.args.recursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssetsUnstructured() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotAssets) != 3 {
				t.Errorf("The number of unstructured asset must be 3 got: %d", len(gotAssets))
				return
			}
			for i := range gotAssets {
				if gotAssets[i].GetKind() != tt.check.kinds[i] {
					t.Errorf("Sort is not correct wanted %s and got: %s", tt.check.kinds[i], gotAssets[i].GetKind())
				}
			}
		})
	}
}

func TestApplier_TemplateAssetsInPathYaml(t *testing.T) {
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	results := make([][]byte, 0)
	for _, y := range assets {
		results = append(results, []byte(y))
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			name:   "success",
			fields: *a,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
			},
			want:    results,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewApplier(tt.fields.reader, tt.fields.values, tt.fields.options)
			if err != nil {
				t.Error(err)
			}
			got, err := a.TemplateAssetsInPathYaml(tt.args.path, tt.args.excluded, tt.args.recursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("Applier.TemplateAssetsInPathYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if len(got) != len(tt.want) {
					t.Errorf("Applier.TemplateAssetsInPathYaml() returns %v yamls, want %v", len(got), len(tt.want))
				}
			}
		})
	}
}

func TestApplier_Assets(t *testing.T) {
	a, err := NewApplier(NewTestReader(), nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	results := make([][]byte, 0)
	for _, y := range assets {
		results = append(results, []byte(y))
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
	}
	tests := []struct {
		name         string
		fields       Applier
		args         args
		wantPayloads [][]byte
		wantErr      bool
	}{
		{
			name:   "success",
			fields: *a,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
			},
			wantPayloads: results,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewApplier(tt.fields.reader, tt.fields.values, tt.fields.options)
			if err != nil {
				t.Error(err)
			}
			gotPayloads, err := a.Assets(tt.args.path, tt.args.excluded, tt.args.recursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("Applier.Assets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPayloads != nil {
				if len(gotPayloads) != len(tt.wantPayloads) {
					t.Errorf("Applier.TemplateAssetsInPathYaml() returns %v yamls, want %v", len(gotPayloads), len(tt.wantPayloads))
				}
			}
		})
	}
}
