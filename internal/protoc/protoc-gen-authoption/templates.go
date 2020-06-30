// Code generated by go-bindata.
// sources:
// templates/auth_method_mapping.go.tmpl
// DO NOT EDIT!

package main

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesAuth_method_mappingGoTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x4f\x6f\xdc\x2c\x10\xc6\xcf\x2f\x9f\x62\x84\x7c\x78\xbb\x4a\x40\xbd\xae\xb4\x87\x2a\x69\xaa\x1e\x92\xb5\xd4\xdc\x23\x62\x26\x18\xad\xf9\x23\x60\xb7\x6d\x10\xdf\xbd\x02\x7b\x1b\x6f\x5b\x55\xe5\x04\xe6\x99\x67\x9e\xf9\x19\xce\xe1\xc6\x49\x04\x85\x16\x83\x48\x28\xe1\xf9\x3b\xf8\xe0\x92\x1b\xae\x15\xda\x6b\x71\x4c\xa3\xc1\x34\x3a\xc9\xe0\x76\x0f\x0f\xfb\x47\xf8\x78\xfb\xf9\x91\x11\xe2\xc5\x70\x10\x0a\x21\x67\x76\xa7\x27\x64\x9f\x5c\x7f\x50\xec\x41\x18\x2c\x85\x10\xa2\x8d\x77\x21\xc1\xff\x04\x00\x80\x2a\xe7\xd4\x84\x4c\xb9\x49\x58\xc5\x5c\x50\x5c\x05\x3f\xd0\x76\x49\xfe\xa3\x4a\xa7\xf1\xf8\xcc\x06\x67\xf8\x20\x5c\xe4\xaf\x3a\x09\x89\x13\xd7\x36\x61\xb0\x62\xe2\xc2\x6b\x5e\xa3\xbc\xd2\x7f\x55\x57\x7f\x1e\x31\x9c\x30\x70\xa3\xa5\x9c\xf0\xab\x08\x48\xc9\x3b\x42\x72\x86\x20\xac\x42\xe8\x22\x6c\x77\x30\xc7\xff\x82\xe1\xa4\x07\x8c\x50\xd3\xf3\xcd\x86\xc0\x06\x72\xee\xe2\x79\x22\xd8\x70\x42\x06\x67\x63\x5a\x7f\x7e\xba\x6f\x6c\xfa\x80\x2f\xfa\x1b\xec\x80\xe6\xdc\xcd\x7e\xfd\x8c\xa7\x14\xb6\x92\x53\x42\x4e\x22\x5c\x18\x7c\x38\xa6\x71\x36\x89\xb0\x83\x36\x23\x9b\xcf\xf7\xc2\x7b\x6d\x15\xe4\x46\xe9\x2d\xb4\xa9\xa1\xbb\xb8\xa8\x6a\xb4\x65\xe5\x0c\x9d\xa9\x7e\x7b\x9f\xaa\xc6\xf9\xa4\x9d\x85\xce\xb0\x7d\xdb\x45\xa0\x95\x17\x5b\x78\xb1\x63\xd2\x53\x64\xa7\xf7\xac\x76\x7d\x9a\xd5\x14\x2e\x0d\xf5\x0b\x08\x2b\x57\xbe\x3f\x77\xac\xc7\x60\x74\x8c\xb5\xc5\xaa\xa6\xfd\x6e\xfe\x77\x0c\x3c\xe7\x33\x90\xed\x32\xf2\x9c\x30\x5f\xd8\xd4\xf5\xd6\x64\xdb\xe0\xfe\xa1\x7b\x29\xf4\xea\xb7\xba\x9b\x11\x87\x43\x2f\x82\x30\xbf\xd4\xb5\x8b\x3b\x8d\x93\x5c\x12\x5c\xd6\x96\xab\xd5\xf4\x68\xcf\x78\x73\x86\xf9\x50\xda\xf3\x41\x2b\xa1\x94\x1f\x01\x00\x00\xff\xff\xf7\x3b\xde\xd5\x3c\x03\x00\x00")

func templatesAuth_method_mappingGoTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesAuth_method_mappingGoTmpl,
		"templates/auth_method_mapping.go.tmpl",
	)
}

func templatesAuth_method_mappingGoTmpl() (*asset, error) {
	bytes, err := templatesAuth_method_mappingGoTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/auth_method_mapping.go.tmpl", size: 828, mode: os.FileMode(420), modTime: time.Unix(1593174056, 0)}
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
	"templates/auth_method_mapping.go.tmpl": templatesAuth_method_mappingGoTmpl,
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
	"templates": &bintree{nil, map[string]*bintree{
		"auth_method_mapping.go.tmpl": &bintree{templatesAuth_method_mappingGoTmpl, map[string]*bintree{}},
	}},
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

