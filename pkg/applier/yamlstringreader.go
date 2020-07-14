package applier

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
)

type YamlStringReader struct {
	yamls []string
}

var _ TemplateReader = &YamlStringReader{
	yamls: []string{""},
}

func (r *YamlStringReader) Asset(
	name string,
) ([]byte, error) {
	i, err := strconv.Atoi(name)
	if err != nil {
		return nil, err
	}
	if i >= len(r.yamls) {
		return nil, fmt.Errorf("Unknown asset %d", i)
	}
	return []byte(r.yamls[i]), nil
}

func (r *YamlStringReader) AssetNames() ([]string, error) {
	keys := make([]string, 0)
	for i := range r.yamls {
		keys = append(keys, strconv.Itoa(i))
	}
	return keys, nil
}

func (*YamlStringReader) ToJSON(
	b []byte,
) ([]byte, error) {
	return yaml.YAMLToJSON(b)
}

func NewYamlStringReader(
	yamls string,
	delimiter string,
) *YamlStringReader {
	yamlsArray := make([]string, 0)
	for _, y := range strings.Split(yamls, delimiter) {
		yamlsArray = append(yamlsArray, strings.TrimSpace(y))
	}

	return &YamlStringReader{
		yamls: yamlsArray,
	}
}
