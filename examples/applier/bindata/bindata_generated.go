// Copyright Contributors to the Open Cluster Management project
// Code generated for package bindata by go-bindata DO NOT EDIT. (@generated)
// sources:
// examples/applier/resources/yamlfilereader/clusterrole.yaml
// examples/applier/resources/yamlfilereader/clusterrolebinding.yaml
// examples/applier/resources/yamlfilereader/namespace.yaml
// examples/applier/resources/yamlfilereader/serviceaccount.yaml
package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _clusterroleYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x90\xb1\x4e\xc3\x40\x0c\x86\xf7\x7b\x0a\x2b\x5d\x69\x10\x1b\xca\x86\x3a\x30\x01\x12\x03\x0b\xea\xe0\xa6\x26\x31\x4d\xce\x87\xed\x6b\x05\x55\xdf\x1d\xa5\x09\x88\xaa\x74\xb2\x4e\xf7\xc9\xfe\xfe\x7f\x06\x0b\x49\x9f\xca\x4d\xeb\xb0\x90\xe8\xca\xab\xec\xa2\x06\x2e\xe0\x2d\xc1\x53\xa2\x08\x8b\x2e\x9b\x93\xc2\x03\x46\x6c\xa8\xa7\xe8\x90\x54\xde\xa9\xf6\x10\x30\xf1\x0b\xa9\xb1\xc4\x0a\x74\x85\x75\x89\xd9\x5b\x51\xfe\x42\x67\x89\xe5\xe6\xd6\x4a\x96\xeb\xed\x4d\xd8\x70\x5c\x57\x3f\xab\x9e\xa5\xa3\xd0\x93\xe3\x1a\x1d\xab\x00\x10\xb1\xa7\x0a\xf6\x7b\x28\xc7\x23\xeb\x09\x7c\xc4\x9e\xe0\x70\x08\x9a\x3b\xb2\x2a\xcc\xe0\xae\xeb\x64\x07\xfd\x08\x01\x36\x83\x8c\x0b\xa8\x38\x3a\x01\xbb\x41\x4d\xea\xfc\xc6\x35\x3a\x85\x39\x60\xe2\x7b\x95\x9c\xac\x82\xd7\xe2\xcf\x97\x4d\x6a\xc5\x32\x00\x28\x99\x64\xad\xe9\x0c\xe2\x26\x72\x6c\x94\x3e\x32\x99\xdb\x91\xdd\x92\xae\x46\x4e\x09\x9d\x8a\x2b\x28\x1a\xf2\x61\x74\x6c\xc7\xb9\x43\xaf\xdb\x62\x79\x59\xb6\x21\x3f\x33\x1b\xe3\x96\x92\x28\xce\xa7\xc7\xbc\xff\xed\xfb\x5f\xd1\x69\xef\x44\xdb\x09\x30\xf4\x76\x84\x2e\x56\x7a\x1a\x66\x88\xb0\x0c\xdf\x01\x00\x00\xff\xff\x18\x5d\xf6\x45\x0e\x02\x00\x00")

func clusterroleYamlBytes() ([]byte, error) {
	return bindataRead(
		_clusterroleYaml,
		"clusterrole.yaml",
	)
}

func clusterroleYaml() (*asset, error) {
	bytes, err := clusterroleYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "clusterrole.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _clusterrolebindingYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\x31\x4f\xc3\x40\x0c\x46\xf7\xfb\x15\x96\x98\x09\x62\x43\xb7\xd1\x0c\x4c\x80\x54\x24\x76\xe7\x62\x12\xd3\xc4\x3e\xf9\x7c\x95\xa0\xea\x7f\x47\xa1\xed\x10\x21\x10\xb3\x9f\xdf\xd3\x77\x05\xad\xe6\x0f\xe3\x61\x74\x68\x55\xdc\xb8\xab\xae\x56\xc0\x15\x7c\x24\x78\xce\x24\xd0\x4e\xb5\x38\x19\x3c\xa2\xe0\x40\x33\x89\x43\x36\x7d\xa7\xe4\x21\x60\xe6\x57\xb2\xc2\x2a\x11\xac\xc3\xd4\x60\xf5\x51\x8d\x3f\xd1\x59\xa5\xd9\xdd\x95\x86\xf5\x66\x7f\x1b\x76\x2c\x7d\xbc\xa8\xb6\x3a\xd1\x86\xa5\x67\x19\xc2\x4c\x8e\x3d\x3a\xc6\x00\x20\x38\x53\x84\xc3\x01\x9a\x53\xab\x3f\xf3\x4f\x38\x13\x1c\x8f\xc1\x74\xa2\x2d\xbd\x2d\x28\x66\x7e\x30\xad\xf9\x8f\x6c\x00\xf8\x51\xfd\x47\xa4\xd4\x6e\xd9\x56\x62\xb8\x3e\xff\xbf\x90\xed\x39\xd1\x7d\x4a\x5a\xc5\x57\x8a\x8d\xaa\x17\x37\xcc\x6b\xe6\xe2\x3a\xa1\x25\x63\xfa\x35\xf9\x7d\x5c\xd8\xaf\x00\x00\x00\xff\xff\x74\x68\x48\xe7\x8c\x01\x00\x00")

func clusterrolebindingYamlBytes() ([]byte, error) {
	return bindataRead(
		_clusterrolebindingYaml,
		"clusterrolebinding.yaml",
	)
}

func clusterrolebindingYaml() (*asset, error) {
	bytes, err := clusterrolebindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "clusterrolebinding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _namespaceYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xcc\xb1\xae\xc2\x30\x0c\x05\xd0\xdd\x5f\x71\xd5\xb7\x3f\x89\x35\x6b\x67\x60\x63\x37\x8d\xd5\x06\x88\x1d\x39\x2e\x12\xaa\xfa\xef\x0c\x20\xf6\xa3\xf3\x87\xd1\xda\xcb\xcb\xbc\x04\x46\xd3\xf0\x72\x5d\xc3\xbc\x23\x0c\xb1\x08\xce\x4d\x14\xe3\x63\xed\x21\x8e\x23\x2b\xcf\x52\x45\x03\xcd\xed\x26\x53\x10\x71\x2b\x17\xf1\x5e\x4c\x13\x9e\x07\xba\x17\xcd\x09\x27\xae\xd2\x1b\x4f\x42\x55\x82\x33\x07\x27\x02\x94\xab\x24\x0c\xdb\x86\xff\xcf\x94\xbf\xf1\x8f\x63\xdf\x07\x7a\x07\x00\x00\xff\xff\x54\x8f\xfb\x9d\x93\x00\x00\x00")

func namespaceYamlBytes() ([]byte, error) {
	return bindataRead(
		_namespaceYaml,
		"namespace.yaml",
	)
}

func namespaceYaml() (*asset, error) {
	bytes, err := namespaceYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "namespace.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _serviceaccountYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x8e\xb1\xce\x82\x40\x10\x84\xfb\x7b\x8a\x0d\x7f\xfd\x9b\xd8\x5e\xa7\xd4\x62\x61\x62\xbf\x1e\x1b\x38\xf5\x76\x2f\x7b\x0b\x09\x21\xbc\xbb\x89\x40\x61\x3b\xf3\x4d\xe6\xfb\x83\x5a\xf2\xa4\xb1\xeb\x0d\x6a\x61\xd3\xf8\x18\x4c\xb4\x80\x09\x58\x4f\x70\xcd\xc4\x50\xbf\x87\x62\xa4\x70\x41\xc6\x8e\x12\xb1\x41\x56\x79\x52\x30\xe7\x30\xc7\x3b\x69\x89\xc2\x1e\xc6\xa3\x7b\x45\x6e\x3d\xdc\x48\xc7\x18\xe8\x14\x82\x0c\x6c\x2e\x91\x61\x8b\x86\xde\x01\x30\x26\xf2\x50\xcd\x33\x1c\xce\x22\x56\x4c\x31\xff\xe2\x0d\x26\x82\x65\xa9\x36\xb8\x64\x0c\xfb\x62\x15\x68\x37\x9f\x66\x6f\xbf\x74\xa1\xa0\x64\xc5\xbb\xff\xed\x23\x4d\x6b\xe4\x3e\x01\x00\x00\xff\xff\x3b\x63\x13\xcc\xe4\x00\x00\x00")

func serviceaccountYamlBytes() ([]byte, error) {
	return bindataRead(
		_serviceaccountYaml,
		"serviceaccount.yaml",
	)
}

func serviceaccountYaml() (*asset, error) {
	bytes, err := serviceaccountYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "serviceaccount.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"clusterrole.yaml":        clusterroleYaml,
	"clusterrolebinding.yaml": clusterrolebindingYaml,
	"namespace.yaml":          namespaceYaml,
	"serviceaccount.yaml":     serviceaccountYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"clusterrole.yaml":        &bintree{clusterroleYaml, map[string]*bintree{}},
	"clusterrolebinding.yaml": &bintree{clusterrolebindingYaml, map[string]*bintree{}},
	"namespace.yaml":          &bintree{namespaceYaml, map[string]*bintree{}},
	"serviceaccount.yaml":     &bintree{serviceaccountYaml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
