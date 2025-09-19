package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Broker struct {
	ID             bson.ObjectID `json:"id,omitzero"              bson:"_id,omitempty"`             // ID
	Name           string        `json:"name"                     bson:"name"`                      // 名字
	Secret         string        `json:"secret"                   bson:"secret,omitempty"`          // 连接密钥
	Status         bool          `json:"status"                   bson:"status,omitempty"`          // 状态
	Goos           string        `json:"goos,omitzero"            bson:"goos,omitempty"`            // GOOS
	Goarch         string        `json:"goarch,omitzero"          bson:"goarch,omitempty"`          // GOARCH
	Protocol       string        `json:"protocol,omitzero"        bson:"protocol,omitempty"`        // 连接协议 tcp/udp
	Config         BrokerConfig  `json:"config,omitempty"         bson:"config,omitempty"`          // 配置
	Networks       NodeNetworks  `json:"networks,omitzero"        bson:"networks,omitempty"`        // 网卡设备
	AliveAt        time.Time     `json:"alive_at,omitzero"        bson:"alive_at,omitempty"`        // 最近心跳时间
	RemoteAddr     string        `json:"remote_addr,omitzero"     bson:"remote_addr,omitempty"`     // 连接的远程地址
	ConnectedAt    time.Time     `json:"connected_at,omitzero"    bson:"connected_at,omitempty"`    // 上线时间
	DisconnectedAt time.Time     `json:"disconnected_at,omitzero" bson:"disconnected_at,omitempty"` // 下线时间
	UpdatedAt      time.Time     `json:"updated_at,omitzero"      bson:"updated_at,omitempty"`      // 数据更新时间
	CreatedAt      time.Time     `json:"created_at,omitzero"      bson:"created_at,omitempty"`      // 数据创建时间
}

type BrokerConfig struct {
	Listen string `json:"listen"`
}
