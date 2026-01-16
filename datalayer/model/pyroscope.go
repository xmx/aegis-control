package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Pyroscope struct {
	ID        bson.ObjectID `bson:"_id,omitempty"        json:"-"`
	Name      string        `bson:"name,omitempty"       json:"name"`
	Address   string        `bson:"address"              json:"address"`
	Username  string        `bson:"username"             json:"username"`
	Password  string        `bson:"password"             json:"password"`
	Enabled   bool          `bson:"enabled"              json:"enabled"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
}
