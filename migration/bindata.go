// Code generated for package migration by go-bindata DO NOT EDIT. (@generated)
// sources:
// migration/000_init_schema.sql
package migration

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

var _migration000_init_schemaSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x93\x4f\x6f\x9b\x40\x10\xc5\xef\x7c\x8a\x77\x04\xb5\x91\xe2\xaa\x95\x2a\x55\x39\x60\xb2\x69\x50\x01\xa7\xb0\x48\xcd\x09\xd6\xde\x09\x59\x95\x1d\x2c\xbc\xa4\xcd\xb7\xaf\x10\x56\x84\x63\xa7\xb2\x94\xb9\xcd\x1f\xfd\x66\x1e\xec\xbb\xb8\xc0\x07\x6b\x9a\x5e\x39\x42\xb9\xf5\xbc\x79\x5e\x38\xe5\xc8\x12\xbb\x25\x35\x86\xbd\x28\x17\xa1\x14\x28\xa2\x5b\x91\x86\x88\x6f\x90\xad\x24\xc4\xaf\xb8\x90\x05\xea\xa7\xce\xe8\x3f\xaa\xfd\x4d\x7d\x8d\x6b\x71\x13\x96\x89\x44\x74\x1b\xe6\x61\x24\x45\x8e\x42\x48\xb4\xca\x19\x5e\x20\x5a\x25\xc9\x88\x99\xd2\xaa\x21\xa6\x5e\xb5\xd5\xc6\x7c\x3b\xbd\x5b\xb0\x3e\xe3\x28\x19\x2e\x13\x81\xda\x3d\x0e\x76\xcd\xca\xb4\xb5\xe7\x7b\x00\x50\x1b\x5d\xe3\x38\xd6\xa6\x31\xec\xfc\x4f\x97\x01\x06\xde\x99\x86\x49\xbf\x1e\x19\xc5\x65\x65\x92\x20\x2c\xe5\xaa\x8a\xb3\x28\x17\xa9\xc8\xe4\xc7\x09\xcb\xca\xd2\x31\xf8\x49\xf5\x9b\x47\xd5\xfb\x8b\x2f\x97\x01\xca\x2c\xfe\x59\x8a\xb7\xb0\x7b\xce\xa6\xb3\xdb\x9e\x76\x3b\xd2\xd5\x0c\x79\x0e\xe7\x24\xc3\x1a\x4b\x95\x7b\xde\x52\x7d\xc8\x38\x1d\x27\x19\x73\x55\xce\xf0\xf3\xf8\x9d\x16\x6f\x11\x8e\xf5\xf4\xa4\x1c\xe9\x4a\xb9\x39\x47\x2b\x47\xce\x58\x82\xa6\x07\x35\xb4\x0e\x51\x99\xe7\x22\x93\x95\x8c\x53\x51\xc8\x30\xbd\x03\x77\x0e\x3c\xb4\xed\x9e\x33\x6c\xf5\xfb\x38\xe8\x18\x13\xe4\x78\x68\xda\x61\x58\xd3\x5f\xf8\xd3\xaf\x0c\x0e\x6b\x33\x19\xaf\x3a\xb3\xc3\xf6\x9d\xbb\x3c\x4e\xc3\xfc\x1e\x3f\xc4\x3d\xfc\xf1\xbd\x05\x5e\x00\x91\x7d\x8f\x33\x81\x2b\xc4\xcc\xdd\xf5\xd2\xc3\x81\x29\x46\x3b\x5c\x61\x70\x0f\x5f\xed\xfa\xb3\x87\x17\x4f\xbc\xd4\xaa\x81\xcd\xa6\xd3\xf4\x7f\x5b\xfc\x0b\x00\x00\xff\xff\x67\x94\x22\xa2\xbc\x03\x00\x00")

func migration000_init_schemaSqlBytes() ([]byte, error) {
	return bindataRead(
		_migration000_init_schemaSql,
		"migration/000_init_schema.sql",
	)
}

func migration000_init_schemaSql() (*asset, error) {
	bytes, err := migration000_init_schemaSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migration/000_init_schema.sql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"migration/000_init_schema.sql": migration000_init_schemaSql,
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
	"migration": &bintree{nil, map[string]*bintree{
		"000_init_schema.sql": &bintree{migration000_init_schemaSql, map[string]*bintree{}},
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
