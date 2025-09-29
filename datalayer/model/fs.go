package model

import (
	"io/fs"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FS struct {
	ID         bson.ObjectID `json:"id"                  bson:"_id,omitempty"`     // ID
	FullPath   string        `json:"full_path"           bson:"full_path"`         // 全路径：home/gopher/main.go
	ParentPath string        `json:"parent_path"         bson:"parent_path"`       // 所在目录
	Name       string        `json:"name"                bson:"name"`              // 名字：main.go
	Size       int64         `json:"size"                bson:"size"`              // 文件大小
	Extension  string        `json:"extension,omitempty" bson:"extension"`         // 转为小写：.go
	FileID     bson.ObjectID `json:"file_id,omitzero"    bson:"file_id,omitempty"` // IsZero 说明是目录
	Checksum   Checksum      `json:"checksum,omitzero"   bson:"checksum,omitempty"`
	CreatedAt  time.Time     `json:"created_at"          bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"          bson:"updated_at"`
}

func (f *FS) IsDir() bool {
	return f.FileID.IsZero()
}

func (f *FS) Mode() fs.FileMode {
	if f.IsDir() {
		return 0o20000000755
	}

	return 0o644
}

type FSs []*FS

// Sort 目录靠前显示，按照文件名字典序排序。
func (fss FSs) Sort() {
	slices.SortFunc(fss, func(a, b *FS) int {
		adir, bdir := a.IsDir(), b.IsDir()
		if adir == bdir {
			return strings.Compare(a.FullPath, b.FullPath)
		}
		if adir {
			return -1
		}
		return 1
	})
}
