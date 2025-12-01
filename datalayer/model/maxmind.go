package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Maxmind struct {
	ID        bson.ObjectID `json:"-"                   bson:"_id,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitzero" bson:"updated_at,omitempty"` // 数据更新时间
	CreatedAt time.Time     `json:"created_at,omitzero" bson:"created_at,omitempty"` // 数据创建时间
}
