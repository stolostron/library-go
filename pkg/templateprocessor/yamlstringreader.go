package templateprocessor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

type YamlStringReader struct {
	Yamls []string
}

var _ TemplateReader = &YamlStringReader{
	Yamls: []string{""},
}

func (r *YamlStringReader) Asset(
	name string,
) ([]byte, error) {
	i, err := strconv.Atoi(name)
	if err != nil {
		return nil, nil
	}
	if i >= len(r.Yamls) {
		return nil, fmt.Errorf("Unknown asset %d", i)
	}
	return []byte(r.Yamls[i]), nil
}

func (r *YamlStringReader) AssetNames() ([]string, error) {
	keys := make([]string, 0)
	for i := range r.Yamls {
		keys = append(keys, strconv.Itoa(i))
	}
	return keys, nil
}

func (*YamlStringReader) ToJSON(
	b []byte,
) ([]byte, error) {
	b, err := yaml.YAMLToJSON(b)
	if err != nil {
		klog.Errorf("err:%s\nyaml:\n%s", err, string(b))
		return nil, err
	}
	return b, nil
}

//NewYamlStringReader returns a YamlStringReader
//yamls: a string of yaml, separeted by the delimiter. Usually "---\n"
//delimiter: the delimiter
func NewYamlStringReader(
	yamls string,
	delimiter string,
) *YamlStringReader {
	yamlsArray := make([]string, 0)
	re := regexp.MustCompile(delimiter)
	ss := re.Split(yamls, -1)
	for _, y := range ss {
		if strings.TrimSpace(y) != "" {
			yamlsArray = append(yamlsArray, strings.TrimSpace(y))
		}
	}

	return &YamlStringReader{
		Yamls: yamlsArray,
	}
}
