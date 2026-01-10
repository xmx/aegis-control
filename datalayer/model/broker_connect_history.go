package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BrokerConnectHistory struct {
	ID         bson.ObjectID     `json:"id"              bson:"_id,omitempty"`
	Broker     bson.ObjectID     `json:"broker_id"       bson:"broker_id"`
	Name       string            `json:"name"            bson:"name"`
	Semver     string            `json:"semver,omitzero" bson:"semver,omitempty"`
	Inet       string            `json:"inet,omitzero"   bson:"inet,omitempty"`
	Goos       string            `json:"goos,omitzero"   bson:"goos,omitempty"`
	Goarch     string            `json:"goarch,omitzero" bson:"goarch,omitempty"`
	TunnelStat TunnelStatHistory `json:"tunnel_stat"     bson:"tunnel_stat"`
}

type BrokerConnectHistories []*BrokerConnectHistory
