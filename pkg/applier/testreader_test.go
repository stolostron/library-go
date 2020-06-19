package applier

import (
	"fmt"

	"github.com/ghodss/yaml"
)

type Test struct{}

var assets = map[string]string{
	"test/clusterrolebinding": `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:test:{{ .ManagedClusterName }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:test:{{ .ManagedClusterName }}
subjects:
- kind: ServiceAccount
  name: {{ .BootstrapServiceAccountName }}
  namespace: {{ .ManagedClusterNamespace }}`,

	"test/serviceaccount": `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: "{{ .BootstrapServiceAccountName }}"
  namespace: "{{ .ManagedClusterNamespace }}"`,

	"test/clusterrole": `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:test:{{ .ManagedClusterName }}
rules:
# Allow managed agent to rotate its certificate
- apiGroups: ["certificates.k8s.io"]
  resources: ["certificatesigningrequests"]
  verbs: ["create", "get", "list", "watch"]
# Allow managed agent to get
- apiGroups: ["cluster.open-cluster-management.io"]
  resources: ["managedclusters"]
  resourceNames: ["{{ .ManagedClusterName }}"]
  verbs: ["get"]`,
}

var values = struct {
	ManagedClusterName          string
	ManagedClusterNamespace     string
	BootstrapServiceAccountName string
}{
	ManagedClusterName:          "mycluster",
	ManagedClusterNamespace:     "myclusterns",
	BootstrapServiceAccountName: "mysa",
}

func (*Test) Asset(name string) ([]byte, error) {
	if s, ok := assets[name]; ok {
		return []byte(s), nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

func (*Test) AssetNames() ([]string, error) {
	keys := make([]string, 0)
	for k := range assets {
		keys = append(keys, k)
	}
	return keys, nil
}

func (*Test) ToJSON(b []byte) ([]byte, error) {
	return yaml.YAMLToJSON(b)
}

func NewTestReader() *Test {
	return &Test{}
}
