// Copyright (c) 2020 Red Hat, Inc.

package applier

import (
	"reflect"
	"testing"
)

func TestYamlFileReader_Asset(t *testing.T) {
	type fields struct {
		rootDirectory string
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
			name: "success",
			fields: fields{
				rootDirectory: "../../test/unit-test/resources/yamlfilereader",
			},
			args: args{
				name: "../../test/unit-test/resources/yamlfilereader/filereader.yaml",
			},
			want: []byte(`apiVersion: fake/v1
kind: Fake
metadata:
  name: {{ .Values }}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YamlFileReader{
				rootDirectory: tt.fields.rootDirectory,
			}
			got, err := y.Asset(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("YamlFileReader.Asset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlFileReader.Asset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYamlFileReader_AssetNames(t *testing.T) {
	type fields struct {
		rootDirectory string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "success",
			fields: fields{
				rootDirectory: "../../test/unit-test/resources/yamlfilereader",
			},
			want: []string{"../../test/unit-test/resources/yamlfilereader/filereader.yaml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &YamlFileReader{
				rootDirectory: tt.fields.rootDirectory,
			}
			if got := r.AssetNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlFileReader.AssetNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
