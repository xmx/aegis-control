package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Firewall HTTP 服务 IP 防火墙配置。
//
// 防火墙可以切换 黑/白名单模式，只能生效一种模式。
// 分为 IP 模式/国家（地区）模式，只能选择一种模式。
type Firewall struct {
	ID           bson.ObjectID `bson:"_id,omitempty"        json:"-"`                      // ID
	Name         string        `bson:"name,omitempty"       json:"name"`                   // 规则名
	Enabled      bool          `bson:"enabled"              json:"enabled"`                // 是否启用该规则（同时最多只能有一条规则生效）
	TrustHeaders []string      `bson:"trust_headers"        json:"trust_headers,omitzero"` // 获取 IP 的可信 Headers（有反代的场景）
	TrustProxies []string      `bson:"trust_proxies"        json:"trust_proxies,omitzero"` // 可信网关（有反代的场景）
	Blacklist    bool          `bson:"blacklist"            json:"blacklist,omitzero"`     // 是否黑名单模式，反之白名单模式
	CountryMode  bool          `bson:"country_mode"         json:"country_mode,omitzero"`  // 是否启用 国家（地区）模式，否则为 IP 模式。
	IPNets       []string      `bson:"ip_nets"              json:"ip_nets,omitzero"`       // IP 模式下的 IP 名单
	Countries    []string      `bson:"countries"            json:"countries,omitzero"`     // 国家（地区）模式下的国家名单，https://www.iso.org/iso-3166-country-codes.html
	CreatedAt    time.Time     `bson:"created_at,omitempty" json:"created_at"`             // 创建时间
	UpdatedAt    time.Time     `bson:"updated_at,omitempty" json:"updated_at"`             // 更新时间
}
