package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Firewall struct {
	ID           bson.ObjectID `bson:"_id,omitempty"           json:"id"`                     // ID
	Name         string        `bson:"name"                    json:"name"`                   // 规则名
	Enabled      bool          `bson:"enabled"                 json:"enabled"`                // 是否启用该规则
	TrustHeaders []string      `bson:"trust_headers,omitempty" json:"trust_headers,omitzero"` // 取 IP 的可信 Headers。
	TrustProxies []string      `bson:"trust_proxies,omitempty" json:"trust_proxies,omitzero"` // 可信网关。
	Blacklist    bool          `bson:"blacklist,omitempty"     json:"blacklist,omitzero"`     // 是否黑名单模式，反之白名单模式。
	Inets        []string      `bson:"inets,omitempty"         json:"inets,omitzero"`         // IP 列表
	Countries    []string      `bson:"countries,omitempty"     json:"countries,omitzero"`     // https://www.iso.org/iso-3166-country-codes.html
	CreatedAt    time.Time     `bson:"created_at,omitempty"    json:"created_at"`             // 创建时间
	UpdatedAt    time.Time     `bson:"updated_at,omitempty"    json:"updated_at"`             // 更新时间
}
