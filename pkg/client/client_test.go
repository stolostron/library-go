// Copyright Contributors to the Open Cluster Management project

package client

import (
	"reflect"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	fakeclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	fakeclientapps "k8s.io/client-go/kubernetes/fake"
)

func TestNewClient(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		url        string
		kubeconfig string
		context    string
		options    client.Options
	}
	tests := []struct {
		name    string
		args    args
		want    client.Client
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: kubeconfigPath,
				context:    "",
				options:    client.Options{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Failed wrongPath",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: "worngPath",
				context:    "",
				options:    client.Options{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.url, tt.args.kubeconfig, tt.args.context, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDefaultClient(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		kubeconfig string
		options    client.Options
	}
	tests := []struct {
		name    string
		args    args
		want    client.Client
		wantErr bool
	}{
		{
			name: "Failed connection refused",
			args: args{
				kubeconfig: kubeconfigPath,
				options:    client.Options{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Failed wrongpaths",
			args: args{
				kubeconfig: "wrongPath",
				options:    client.Options{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDefaultClient(tt.args.kubeconfig, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKubeClient(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		url        string
		kubeconfig string
		context    string
	}
	tests := []struct {
		name    string
		args    args
		want    kubernetes.Interface
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: kubeconfigPath,
				context:    "",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Failed wrong path",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: "wrongPath",
				context:    "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKubeClient(tt.args.url, tt.args.kubeconfig, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubeClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//we can not test the content
		})
	}
}

func TestNewDefaultKubeClient(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name    string
		args    args
		want    kubernetes.Interface
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				kubeconfig: kubeconfigPath,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Failed wrong path",
			args: args{
				kubeconfig: "wrongPath",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDefaultKubeClient(tt.args.kubeconfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultKubeClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewKubeClientDynamic(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		url        string
		kubeconfig string
		context    string
	}
	tests := []struct {
		name    string
		args    args
		want    dynamic.Interface
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: kubeconfigPath,
				context:    "",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Failed wrong path",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: "wrongPath",
				context:    "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKubeClientDynamic(tt.args.url, tt.args.kubeconfig, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubeClientDynamic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewDefaultKubeClientDynamic(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name    string
		args    args
		want    dynamic.Interface
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				kubeconfig: kubeconfigPath,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Failed wrong path",
			args: args{
				kubeconfig: "wrongPath",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDefaultKubeClientDynamic(tt.args.kubeconfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultKubeClientDynamic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewKubeClientAPIExtension(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		url        string
		kubeconfig string
		context    string
	}
	tests := []struct {
		name    string
		args    args
		want    clientset.Interface
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: kubeconfigPath,
				context:    "",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Failed wrong path",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: "wrongPath",
				context:    "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKubeClientAPIExtension(tt.args.url, tt.args.kubeconfig, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubeClientAPIExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewDefaultKubeClientAPIExtension(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name    string
		args    args
		want    clientset.Interface
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				kubeconfig: kubeconfigPath,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Failed wrong path",
			args: args{
				kubeconfig: "wrongPath",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDefaultKubeClientAPIExtension(tt.args.kubeconfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultKubeClientAPIExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

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

func TestHaveDeploymentsInNamespace(t *testing.T) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
	}

	deploymentMinAvailable := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{Reason: "MinimumReplicasAvailable", Status: corev1.ConditionTrue},
			},
		},
	}
	deploymentNoMinAvailable := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mydeployment",
			Namespace: "mynamespace",
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{Reason: "MinimumReplicasAvailable", Status: corev1.ConditionFalse},
			},
		},
	}

	client := fakeclientapps.NewSimpleClientset(deployment)
	clientMinAvailable := fakeclientapps.NewSimpleClientset(deploymentMinAvailable)
	clientNoMinAvailable := fakeclientapps.NewSimpleClientset(deploymentNoMinAvailable)

	type args struct {
		client                  kubernetes.Interface
		namespace               string
		expectedDeploymentNames []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "deployment exists",
			args: args{
				client:                  client,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantErr: false,
		},
		{
			name: "all deployment not present",
			args: args{
				client:                  client,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment", "notexists"},
			},
			wantErr: true,
		},
		{
			name: "Deployment no minimum available",
			args: args{
				client:                  clientNoMinAvailable,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantErr: true,
		},
		{
			name: "Deployment minimum available",
			args: args{
				client:                  clientMinAvailable,
				namespace:               "mynamespace",
				expectedDeploymentNames: []string{"mydeployment"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HaveDeploymentsInNamespace(tt.args.client, tt.args.namespace, tt.args.expectedDeploymentNames); (err != nil) != tt.wantErr {
				t.Errorf("HaveDeploymentsInNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
