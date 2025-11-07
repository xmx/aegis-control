package model

import (
	"runtime/debug"
	"unsafe"
)

type BuildInfo struct {
	GoVersion string         `bson:"go_version,omitempty" json:"go_version,omitzero"`
	Path      string         `bson:"path,omitempty"       json:"path,omitzero"`
	Main      BuildModule    `bson:"main"                 json:"main"`
	Deps      []*BuildModule `bson:"deps,omitempty"       json:"deps,omitzero"`
	Settings  []BuildSetting `bson:"settings,omitempty"   json:"settings,omitzero"`
}

type BuildModule struct {
	Path    string       `bson:"path,omitempty"    json:"path,omitzero"`    // module path
	Version string       `bson:"version,omitempty" json:"version,omitzero"` // module version
	Sum     string       `bson:"sum,omitempty"     json:"sum,omitzero"`     // checksum
	Replace *BuildModule `bson:"replace,omitempty" json:"replace,omitzero"` // replaced by this module
}

type BuildSetting struct {
	Key   string `bson:"key,omitempty"   json:"key,omitzero"`
	Value string `bson:"value,omitempty" json:"value,omitzero"`
}

func FormatBuildInfo(bi *debug.BuildInfo) *BuildInfo {
	return (*BuildInfo)(unsafe.Pointer(bi))
}
