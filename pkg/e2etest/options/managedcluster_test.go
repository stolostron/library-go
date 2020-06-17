package options

import (
	"reflect"
	"testing"
)

func TestGetManagedClusterKubeConfigs(t *testing.T) {
	type args struct {
		configDir string
		scenario  string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				configDir: "../../../test/unit-test/resources/clusters",
				scenario:  "fake-scenario1",
			},
			want: map[string]string{
				"cluster1": "../../../test/unit-test/resources/clusters/fake-scenario1/cluster1/kubeconfig.yaml",
				"cluster2": "../../../test/unit-test/resources/clusters/fake-scenario1/cluster2/kubeconfig.yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetManagedClusterKubeConfigs(tt.args.configDir, tt.args.scenario)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetManagedClusterKubeConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManagedClusterKubeConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}
