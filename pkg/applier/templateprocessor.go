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
	AssetNames() []string
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
//
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

//TemplateAssets render the given templates with the provided config
//The assets are not sorted
func (a *TemplateProcessor) TemplateAssets(templateNames []string, values interface{}) ([][]byte, error) {
	results := make([][]byte, len(templateNames))
	for i, templateName := range templateNames {
		result, err := a.TemplateAsset(templateName, values)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

//TemplateAsset render the given template with the provided config
func (a *TemplateProcessor) TemplateAsset(templateName string, values interface{}) ([]byte, error) {
	var buf bytes.Buffer
	b, err := a.reader.Asset(templateName)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(templateName).Parse(string(b))
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
func (a *TemplateProcessor) TemplateAssetsInPathYaml(path string, excluded []string, recursive bool, values interface{}) ([][]byte, error) {
	us, err := a.TemplateAssetsInPathUnstructured(path, excluded, recursive, values)
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
func (a *TemplateProcessor) AssetNamesInPath(path string, excluded []string, recursive bool) []string {
	results := make([]string, 0)
	names := a.reader.AssetNames()
	for _, name := range names {
		if isExcluded(name, excluded) {
			continue
		}
		if (recursive && strings.HasPrefix(name, path)) ||
			(!recursive && filepath.Dir(name) == path) {
			results = append(results, name)
		}
	}
	return results
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
func (a *TemplateProcessor) Assets(path string, excluded []string, recursive bool) (payloads [][]byte, err error) {
	names := a.AssetNamesInPath(path, excluded, recursive)

	for _, name := range names {
		b, err := a.reader.Asset(name)
		if err != nil {
			return nil, err
		}
		payloads = append(payloads, b)
	}
	return payloads, nil
}

// TemplateAssetsInPathUnstructured returns all assets in a []unstructured.Unstructured and sort them
// The []unstructured.Unstructured are sorted following the order defined in variable kindsOrder
func (a *TemplateProcessor) TemplateAssetsInPathUnstructured(
	path string,
	excluded []string,
	recursive bool,
	values interface{}) (assets []*unstructured.Unstructured, err error) {
	templateNames := a.AssetNamesInPath(path, excluded, recursive)
	templatedAssets, err := a.TemplateAssets(templateNames, values)
	if err != nil {
		return nil, err
	}
	assets = make([]*unstructured.Unstructured, len(templateNames))

	for i, b := range templatedAssets {
		j, err := a.reader.ToJSON(b)
		if err != nil {
			return nil, err
		}
		u := &unstructured.Unstructured{}
		err = u.UnmarshalJSON(j)
		if err != nil {
			return nil, err
		}
		assets[i] = u
	}
	a.sortUnstructuredForApply(assets)
	return assets, nil
}

//sortUnstructuredForApply sorts a list on unstructured
func (a *TemplateProcessor) sortUnstructuredForApply(objects []*unstructured.Unstructured) {
	sort.Slice(objects[:], func(i, j int) bool {
		return a.less(objects[i], objects[j])
	})
}

func (a *TemplateProcessor) less(u1, u2 *unstructured.Unstructured) bool {
	if a.weight(u1) == a.weight(u2) {
		if u1.GetNamespace() == u2.GetNamespace() {
			return u1.GetName() < u2.GetName()
		}
		return u1.GetNamespace() < u2.GetNamespace()
	}
	return a.weight(u1) < a.weight(u2)
}

func (a *TemplateProcessor) weight(u *unstructured.Unstructured) int {
	kind := u.GetKind()
	for i, k := range a.options.KindsOrder {
		if k == kind {
			return i
		}
	}
	return len(a.options.KindsOrder)
}
