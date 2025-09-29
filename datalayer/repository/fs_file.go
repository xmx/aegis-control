package repository

import (
	"errors"
	"io/fs"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type fsFile struct {
	stm  *mongo.GridFSDownloadStream
	file *model.FS
}

func (f *fsFile) Stat() (fs.FileInfo, error) {
	return &fsFileInfo{file: f.file}, nil
}

func (f *fsFile) Read(b []byte) (int, error) {
	if f.stm == nil || f.file.IsDir() {
		name := f.file.Name
		err := errors.New("is a directory")
		return 0, &fs.PathError{Op: "read", Path: name, Err: err}
	}

	return f.stm.Read(b)
}

func (f *fsFile) Close() error {
	if stm := f.stm; stm != nil {
		return stm.Close()
	}

	return nil
}

type fsFileInfo struct {
	file *model.FS
}

func (f *fsFileInfo) Name() string {
	return f.file.Name
}

func (f *fsFileInfo) Size() int64 {
	return f.file.Size
}

func (f *fsFileInfo) Mode() fs.FileMode {
	return f.file.Mode()
}

func (f *fsFileInfo) ModTime() time.Time {
	return f.file.UpdatedAt
}

func (f *fsFileInfo) IsDir() bool {
	return f.file.IsDir()
}

func (f *fsFileInfo) Sys() any {
	return f.file
}

type fsDirEntry struct {
	file *model.FS
}

func (f *fsDirEntry) Name() string {
	return f.file.Name
}

func (f *fsDirEntry) IsDir() bool {
	return f.file.IsDir()
}

func (f *fsDirEntry) Type() fs.FileMode {
	return f.file.Mode()
}

func (f *fsDirEntry) Info() (fs.FileInfo, error) {
	return &fsFileInfo{file: f.file}, nil
}
