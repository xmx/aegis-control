package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VictoriaMetrics struct {
	ID        bson.ObjectID `bson:"_id,omitempty"        json:"-"`
	Name      string        `bson:"name,omitempty"       json:"name"`
	Method    string        `bson:"method"               json:"method"`
	Address   string        `bson:"address"              json:"address"`
	Header    HTTPHeader    `bson:"header"               json:"header"`
	Enabled   bool          `bson:"enabled"              json:"enabled"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
}
