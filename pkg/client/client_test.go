package client

import (
	"reflect"
	"testing"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestHaveCRDs(t *testing.T) {
	// This test is covered by
	// github.com/open-cluster-management/library-go/pkg/apis/meta/v1/crd.TestHaveCRDs(t)
}

func TestHaveDeploymentsInNamespace(t *testing.T) {
	// This test is covered by
	// github.com/open-cluster-management/library-go/pkg/apis/meta/v1/deployment.TestHaveDeploymentsInNamespace(t)")
}

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
