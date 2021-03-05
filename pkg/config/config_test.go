// Copyright Contributors to the Open Cluster Management project

package config

import (
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"testing"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestLoadConfig(t *testing.T) {
	kubeconfigPath := "../../test/unit/resources/config/kubeconfig.yaml"
	apiConfig, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		t.Error(err)
	}

	userconfigexists := true
	if usr, err := user.Current(); err == nil {
		if _, err := os.Stat(filepath.Join(usr.HomeDir, ".kube", "config")); os.IsNotExist(err) {
			userconfigexists = false
		}
	}

	inCluster := false
	if _, err = rest.InClusterConfig(); err == nil {
		inCluster = true
	}

	config, err := clientcmd.NewDefaultClientConfig(
		*apiConfig,
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		t.Error(err)
	}
	type args struct {
		url        string
		kubeconfig string
		context    string
	}
	tests := []struct {
		name    string
		args    args
		want    *rest.Config
		wantErr bool
	}{
		{
			name: "success from file current context",
			args: args{
				url:        "",
				kubeconfig: kubeconfigPath,
				context:    "",
			},
			want:    config,
			wantErr: false,
		},
		{
			name: "success from file specified context",
			args: args{
				url:        "",
				kubeconfig: kubeconfigPath,
				context:    "kind-rcm-e2e-test",
			},
			want:    config,
			wantErr: false,
		},
		{
			name: "success from file specified url",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: kubeconfigPath,
				context:    "",
			},
			want:    config,
			wantErr: false,
		},
		{
			name: "failed wrong path",
			args: args{
				url:        "https://127.0.0.1:32878",
				kubeconfig: "wrong-path",
				context:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Succeed user or InCluster",
			args: args{
				url:        "",
				kubeconfig: "",
				context:    "",
			},
			want:    nil,
			wantErr: !(userconfigexists || inCluster),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test name: %s", tt.name)
			got, err := LoadConfig(tt.args.url, tt.args.kubeconfig, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
