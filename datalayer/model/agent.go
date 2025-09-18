package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Agent struct {
	ID             bson.ObjectID         `json:"id,omitzero"              bson:"_id,omitempty"`             // ID
	MachineID      string                `json:"machine_id"               bson:"machine_id"`                // 机器码（全局唯一）
	Goos           string                `json:"goos"                     bson:"goos"`                      // GOOS
	Goarch         string                `json:"goarch"                   bson:"goarch"`                    // GOARCH
	Status         bool                  `json:"status"                   bson:"status"`                    // 节点状态
	Broker         *AgentConnectedBroker `json:"broker,omitzero"          bson:"broker,omitempty"`          // agent 所在的 broker
	AliveAt        time.Time             `json:"alive_at,omitzero"        bson:"alive_at,omitempty"`        // 最近心跳时间
	Protocol       string                `json:"protocol,omitzero"        bson:"protocol,omitempty"`        // 连接协议 tcp/udp
	RemoteAddr     string                `json:"remote_addr,omitzero"     bson:"remote_addr,omitempty"`     // 连接的远程地址
	ConnectedAt    time.Time             `json:"connected_at,omitzero"    bson:"connected_at,omitempty"`    // 上线时间
	DisconnectedAt time.Time             `json:"disconnected_at,omitzero" bson:"disconnected_at,omitempty"` // 下线时间
	CreatedAt      time.Time             `json:"created_at"               bson:"created_at,omitempty"`      // 创建时间
	UpdatedAt      time.Time             `json:"updated_at"               bson:"updated_at,omitempty"`      // 更新时间
}

type AgentConnectedBroker struct {
	ID   bson.ObjectID `json:"id,omitzero"   bson:"id,omitempty"`
	Name string        `json:"name,omitzero" bson:"name,omitempty"`
}
