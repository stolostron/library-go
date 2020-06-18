// Copyright (c) 2020 Red Hat, Inc.

package applier

import (
	"bytes"
	goerr "errors"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//TemplateProcessor this structure holds all object for the applier
type TemplateProcessor struct {
	//reader a TemplateReader to read the data source
	reader TemplateReader
	//Options to configure the applier
	options *Options
}

//TemplateReader defines the needed functions
type TemplateReader interface {
	//Retreive an asset from the data source
	Asset(templatePath string) ([]byte, error)
	//List all available assets in the data source
	AssetNames() ([]string, error)
	//Transform the assets into a JSON. This is used to transform the asset into an unstructrued.Unstructured object.
	//For example: if the asset is a yaml, you can use yaml.YAMLToJSON(b []byte) as implementation as it is shown in
	//testread_test.go
	ToJSON(b []byte) ([]byte, error)
}

//Options defines for the available options for the applier
type Options struct {
	//Override the default order, it contains the kind order which the applier must use before applying all resources.
	KindsOrder []string
}

var log = logf.Log.WithName("applier")

//defaultKindsOrder the default order
var defaultKindsOrder = []string{
	"CustomResourceDefinition",
	"ClusterRole",
	"ClusterRoleBinding",
	"Namespace",
	"Secret",
	"ServiceAccount",
	"Role",
	"RoleBinding",
	"ConfigMap",
	"Deployment",
}

//NewTemplateProcessor creates a new applier
//reader: The TemplateReader to use to read the templates
//options: The possible options for the templateprocessor
func NewTemplateProcessor(
	reader TemplateReader,
	options *Options,
) (*TemplateProcessor, error) {
	if reader == nil {
		return nil, goerr.New("reader is nil")
	}
	if options == nil {
		options = &Options{}
	}
	if options.KindsOrder == nil {
		options.KindsOrder = defaultKindsOrder
	}
	return &TemplateProcessor{
		reader:  reader,
		options: options,
	}, nil
}

//TemplateAssets render the given templates with the provided values
//The assets are not sorted
func (tp *TemplateProcessor) TemplateAssets(templateNames []string, values interface{}) ([][]byte, error) {
	results := make([][]byte, len(templateNames))
	for i, templateName := range templateNames {
		result, err := tp.TemplateAsset(templateName, values)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

//TemplateAsset render the given template with the provided values
func (tp *TemplateProcessor) TemplateAsset(templateName string, values interface{}) ([]byte, error) {
	b, err := tp.reader.Asset(templateName)
	if err != nil {
		return nil, err
	}
	return tp.TemplateBytes(b, values)
}

//TemplateBytes render the given template with the provided values
func (tp *TemplateProcessor) TemplateBytes(b []byte, values interface{}) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("yamls").Parse(string(b))
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// TemplateAssetsInPathYaml returns all assets in a path using the provided config.
// The assets are sorted following the order defined in variable kindsOrder
func (tp *TemplateProcessor) TemplateAssetsInPathYaml(path string, excluded []string, recursive bool, values interface{}) ([][]byte, error) {
	us, err := tp.TemplateAssetsInPathUnstructured(path, excluded, recursive, values)
	if err != nil {
		return nil, err
	}

	results := make([][]byte, len(us))

	for i, u := range us {
		j, err := u.MarshalJSON()
		if err != nil {
			return nil, err
		}
		y, err := yaml.JSONToYAML(j)
		if err != nil {
			return nil, err
		}
		results[i] = y
	}
	return results, nil
}

//AssetNamesInPath returns all asset names with a given path and
// subpath if recursive is set to true, it excludes the assets contained in the excluded parameter
func (tp *TemplateProcessor) AssetNamesInPath(path string, excluded []string, recursive bool) ([]string, error) {
	results := make([]string, 0)
	names, err := tp.reader.AssetNames()
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if isExcluded(name, excluded) {
			continue
		}
		if (recursive && strings.HasPrefix(name, path)) ||
			(!recursive && filepath.Dir(name) == path) {
			results = append(results, name)
		}
	}
	return results, nil
}

func isExcluded(name string, excluded []string) bool {
	if excluded == nil {
		return false
	}
	for _, e := range excluded {
		if e == name {
			return true
		}
	}
	return false
}

//Assets returns all assets with a given path and
// subpath if recursive set to true, it excludes the assets contained in the excluded parameter
func (tp *TemplateProcessor) Assets(path string, excluded []string, recursive bool) (payloads [][]byte, err error) {
	names, err := tp.AssetNamesInPath(path, excluded, recursive)
	if err != nil {
		return nil, err
	}

	for _, name := range names {
		b, err := tp.reader.Asset(name)
		if err != nil {
			return nil, err
		}
		payloads = append(payloads, b)
	}
	return payloads, nil
}

// TemplateAssetsInPathUnstructured returns all assets in a []unstructured.Unstructured and sort them
// The []unstructured.Unstructured are sorted following the order defined in variable kindsOrder
func (tp *TemplateProcessor) TemplateAssetsInPathUnstructured(
	path string,
	excluded []string,
	recursive bool,
	values interface{}) (assets []*unstructured.Unstructured, err error) {
	templateNames, err := tp.AssetNamesInPath(path, excluded, recursive)
	if err != nil {
		return nil, err
	}
	templatedAssets, err := tp.TemplateAssets(templateNames, values)
	if err != nil {
		return nil, err
	}
	assets, err = tp.BytesArrayToUnstructured(templatedAssets)
	if err != nil {
		return nil, err
	}
	tp.sortUnstructuredForApply(assets)
	return assets, nil
}

// TemplateBytesUnstructured returns all assets defined in a []byte (separted by the delimiter)
// in a []unstructured.Unstructured and sort them
// The []unstructured.Unstructured are sorted following the order defined in variable kindsOrder
func (tp *TemplateProcessor) TemplateBytesUnstructured(
	assets []byte,
	values interface{},
	delimiter string) (us []*unstructured.Unstructured, err error) {
	templatedAssets, err := tp.TemplateBytes(assets, values)
	if err != nil {
		return nil, err
	}
	assetsString := string(templatedAssets)
	ys := strings.Split(assetsString, delimiter)
	templatedAssetsArray := make([][]byte, 0)
	for _, y := range ys {
		if len(strings.TrimSpace(y)) == 0 {
			continue
		}
		templatedAssetsArray = append(templatedAssetsArray, []byte(y))
	}
	us, err = tp.BytesArrayToUnstructured(templatedAssetsArray)
	if err != nil {
		return nil, err
	}
	tp.sortUnstructuredForApply(us)
	return us, nil

}

//BytesArrayToUnstructured transform a [][]byte to an []*unstructured.Unstructured using the TemplateProcessor reader
func (tp *TemplateProcessor) BytesArrayToUnstructured(assets [][]byte) (us []*unstructured.Unstructured, err error) {
	us = make([]*unstructured.Unstructured, len(assets))
	for i, b := range assets {
		u, err := tp.BytesToUnstructured(b)
		if err != nil {
			return nil, err
		}
		us[i] = u
	}
	return us, nil
}

//BytesToUnstructured transform a []byte to an *unstructured.Unstructured using the TemplateProcessor reader
func (tp *TemplateProcessor) BytesToUnstructured(asset []byte) (*unstructured.Unstructured, error) {
	j, err := tp.reader.ToJSON(asset)
	if err != nil {
		return nil, err
	}
	u := &unstructured.Unstructured{}
	err = u.UnmarshalJSON(j)
	if err != nil {
		return nil, err
	}
	return u, nil
}

//sortUnstructuredForApply sorts a list on unstructured
func (tp *TemplateProcessor) sortUnstructuredForApply(us []*unstructured.Unstructured) {
	sort.Slice(us[:], func(i, j int) bool {
		return tp.less(us[i], us[j])
	})
}

func (tp *TemplateProcessor) less(u1, u2 *unstructured.Unstructured) bool {
	if tp.weight(u1) == tp.weight(u2) {
		if u1.GetNamespace() == u2.GetNamespace() {
			return u1.GetName() < u2.GetName()
		}
		return u1.GetNamespace() < u2.GetNamespace()
	}
	return tp.weight(u1) < tp.weight(u2)
}

func (tp *TemplateProcessor) weight(u *unstructured.Unstructured) int {
	kind := u.GetKind()
	for i, k := range tp.options.KindsOrder {
		if k == kind {
			return i
		}
	}
	return len(tp.options.KindsOrder)
}
