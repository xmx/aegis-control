package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Setting struct {
	ID        bson.ObjectID   `json:"-"                   bson:"_id,omitempty"`
	Exposes   ExposeAddresses `json:"exposes"             bson:"exposes"`              // 服务暴露地址
	UpdatedAt time.Time       `json:"updated_at,omitzero" bson:"updated_at,omitempty"` // 数据更新时间
	CreatedAt time.Time       `json:"created_at,omitzero" bson:"created_at,omitempty"` // 数据创建时间
}

type SettingData struct {
	Exposes ExposeAddresses `json:"exposes"`
}

type ExposeAddress struct {
	Name string `json:"name" bson:"name"`
	Addr string `json:"addr" bson:"addr"`
}

type ExposeAddresses []*ExposeAddress

func (eas ExposeAddresses) Addresses() []string {
	rets := make([]string, 0, 10)
	uniq := make(map[string]struct{}, 8)
	for _, ea := range eas {
		addr := ea.Addr
		if _, exists := uniq[addr]; !exists {
			uniq[addr] = struct{}{}
			rets = append(rets, addr)
		}
	}

	return rets
}
