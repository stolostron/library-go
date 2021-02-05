package templateprocessor

import (
	"bytes"
	"encoding/base64"
	"text/template"

	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

var tmpl *template.Template

func ApplierFuncMap(t *template.Template) template.FuncMap {
	tmpl = t
	return template.FuncMap(GenericFuncMap())
}

// GenericFuncMap returns a copy of the basic function map as a map[string]interface{}.
func GenericFuncMap() map[string]interface{} {
	gfm := make(map[string]interface{}, len(genericMap))
	for k, v := range genericMap {
		gfm[k] = v
	}
	return gfm
}

var genericMap = map[string]interface{}{
	"toYaml":       toYaml,
	"encodeBase64": encodeBase64,
	"include":      include,
}

func toYaml(o interface{}) (string, error) {
	m, err := yaml.Marshal(o)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	klog.V(5).Infof(string(m))
	return string(m), nil
}

func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func include(name string, data interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
