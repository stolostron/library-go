package applier

import (
	"reflect"
	"testing"
)

var assetsYamls = []string{`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ .ManagedClusterName }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: {{ .ManagedClusterName }}
subjects:
- kind: ServiceAccount
  name: {{ .BootstrapServiceAccountName }}
  namespace: {{ .ManagedClusterNamespace }}`, `apiVersion: v1
kind: ServiceAccount
metadata:
  name: "{{ .BootstrapServiceAccountName }}"
  namespace: "{{ .ManagedClusterNamespace }}"
secrets:
- name: mysecret`, `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ .ManagedClusterName }}
rules:
# Allow managed agent to rotate its certificate
- apiGroups: ["certificates.k8s.io"]
  resources: ["certificatesigningrequests"]
  verbs: ["create", "get", "list", "watch"]
# Allow managed agent to get
- apiGroups: ["cluster.open-cluster-management.io"]
  resources: ["managedclusters"]
  resourceNames: ["{{ .ManagedClusterName }}"]
  verbs: ["get"]`}

func TestNewYamlStringReader(t *testing.T) {

	type args struct {
		yamls     string
		delimiter string
	}
	tests := []struct {
		name string
		args args
		want *YamlStringReader
	}{
		{
			name: "create",
			args: args{
				yamls:     assetsYaml,
				delimiter: "---",
			},
			want: &YamlStringReader{
				yamls: assetsYamls,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewYamlStringReader(tt.args.yamls, tt.args.delimiter)
			t.Log(len(got.yamls))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewYamlStringReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYamlStringReader_Asset(t *testing.T) {
	type fields struct {
		yamls []string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "get",
			fields: fields{
				yamls: assetsYamls,
			},
			args: args{
				name: "1",
			},
			want:    []byte(assetsYamls[1]),
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				yamls: assetsYamls,
			},
			args: args{
				name: "hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-exist",
			fields: fields{
				yamls: assetsYamls,
			},
			args: args{
				name: "3",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &YamlStringReader{
				yamls: tt.fields.yamls,
			}
			got, err := r.Asset(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("YamlStringReader.Asset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlStringReader.Asset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYamlStringReader_AssetNames(t *testing.T) {
	type fields struct {
		yamls []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "full list",
			fields: fields{
				yamls: assetsYamls,
			},
			want:    []string{"0", "1", "2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &YamlStringReader{
				yamls: tt.fields.yamls,
			}
			got, err := r.AssetNames()
			if (err != nil) != tt.wantErr {
				t.Errorf("YamlStringReader.AssetNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlStringReader.AssetNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
