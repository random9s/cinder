package sys

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/random9s/cinder/sys/mem"
)

//FileSystem represents an in memory file system
type FileSystem map[string]*zip.File

//Unzip a file into a FileSystem type
func Unzip(file *mem.File) (FileSystem, error) {
	//Create new zip file reader
	zr, err := zip.NewReader(file, file.Size())
	if err != nil {
		return nil, fmt.Errorf("sys.Unzip: %s", err)
	}

	// Iterate over files in archive
	var fs = make(map[string]*zip.File)
	for _, f := range zr.File {
		fs[f.Name] = f
	}

	return fs, nil
}

//FilesWithPrefix return all files with the given prefix
func (fs FileSystem) FilesWithPrefix(prefix string) []*zip.File {
	var files = make([]*zip.File, 0)
	for k, v := range fs {
		if strings.HasPrefix(k, prefix) {
			files = append(files, v)
		}
	}
	return files
}

//OpenFile copies uncompressed contents into memfile
func (fs FileSystem) OpenFile(name string) (*mem.File, error) {
	v, ok := fs[name]
	if !ok {
		return nil, fmt.Errorf("sys.OpenFile: file \"%s\" could not be located", name)
	}

	fp, err := v.Open()
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	b, err := ioutil.ReadAll(fp)
	return mem.NewFile(b), err
}

//Remove removes element from map
func (fs FileSystem) Remove(name string) {
	delete(fs, name)
}
