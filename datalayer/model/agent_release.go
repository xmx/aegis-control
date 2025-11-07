package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AgentRelease struct {
	ID        bson.ObjectID `bson:"_id,omitempty"        json:"id"`
	FileID    bson.ObjectID `bson:"file_id,omitempty"    json:"-"`
	Filename  string        `bson:"filename,omitempty"   json:"filename"`
	Goos      string        `bson:"goos,omitempty"       json:"goos"`
	Goarch    string        `bson:"goarch,omitempty"     json:"goarch"`
	Length    int64         `bson:"length,omitempty"     json:"length"`
	Semver    string        `bson:"semver"               json:"semver"`
	Version   uint64        `bson:"version,omitempty"    json:"version"`
	BuildInfo *BuildInfo    `bson:"build_info,omitempty" json:"build_info"`
	Changelog string        `bson:"changelog,omitempty"  json:"changelog"`
	Checksum  Checksum      `bson:"checksum,omitempty"   json:"checksum"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at"`
}

type AgentReleases []*AgentRelease
