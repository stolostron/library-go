package applier

var values = struct {
	ManagedClusterName          string
	ManagedClusterNamespace     string
	BootstrapServiceAccountName string
}{
	ManagedClusterName:          "mycluster",
	ManagedClusterNamespace:     "myclusterns",
	BootstrapServiceAccountName: "mysa",
}

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
