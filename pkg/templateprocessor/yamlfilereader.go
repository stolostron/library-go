package templateprocessor

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

type YamlFileReader struct {
	path string
}

var _ TemplateReader = &YamlFileReader{
	path: "",
}

func (r *YamlFileReader) Asset(
	name string,
) ([]byte, error) {
	return ioutil.ReadFile(filepath.Clean(filepath.Join(r.path, name)))
}

func (r *YamlFileReader) AssetNames() ([]string, error) {
	keys := make([]string, 0)
	fi, err := os.Stat(r.path)
	if err != nil {
		return keys, err
	}
	if fi.Mode().IsDir() {
		err = filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
			if info != nil {
				if !info.IsDir() {
					newPath, err := filepath.Rel(r.path, path)
					if err != nil {
						return err
					}
					keys = append(keys, newPath)
				}
			}
			return nil
		})
	} else {
		keys = append(keys, "")
	}
	return keys, err
}

func (*YamlFileReader) ToJSON(
	b []byte,
) ([]byte, error) {
	b, err := yaml.YAMLToJSON(b)
	if err != nil {
		klog.Errorf("err:%s\nyaml:\n%s", err, string(b))
		return nil, err
	}
	return b, nil
}

func NewYamlFileReader(
	path string,
) *YamlFileReader {
	return &YamlFileReader{
		path: path,
	}
}
