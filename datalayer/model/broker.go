package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Broker struct {
	ID          bson.ObjectID `json:"id,omitzero"              bson:"_id,omitempty"`          // ID
	Name        string        `json:"name"                     bson:"name"`                   // 名字
	Secret      string        `json:"secret"                   bson:"secret,omitempty"`       // 连接密钥
	Status      bool          `json:"status"                   bson:"status,omitempty"`       // 状态
	Config      BrokerConfig  `json:"config,omitempty"         bson:"config,omitempty"`       // 配置
	Networks    NodeNetworks  `json:"networks,omitzero"        bson:"networks,omitempty"`     // 网卡设备
	TunnelStat  *TunnelStat   `json:"tunnel_stat,omitzero"     bson:"tunnel_stat,omitempty"`  // 通道连接状态
	ExecuteStat *ExecuteStat  `json:"execute_stat,omitzero"    bson:"execute_stat,omitempty"` // broker 执行程序信息
	UpdatedAt   time.Time     `json:"updated_at,omitzero"      bson:"updated_at,omitempty"`   // 数据更新时间
	CreatedAt   time.Time     `json:"created_at,omitzero"      bson:"created_at,omitempty"`   // 数据创建时间
}

type BrokerConfig struct {
	Server BrokerServerConfig `json:"server" bson:"server"`
	Logger BrokerLoggerConfig `json:"logger" bson:"logger"`
}

type BrokerServerConfig struct {
	Addr string `json:"addr" bson:"addr"`
}

type BrokerLoggerConfig struct {
	Level      string `json:"level"      bson:"level"      validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
	Console    bool   `json:"console"    bson:"console"`
	Filename   string `json:"filename"   bson:"filename"`
	MaxSize    int    `json:"maxsize"    bson:"maxsize"    validate:"gte=0"`
	MaxAge     int    `json:"maxage"     bson:"maxage"     validate:"gte=0"`
	MaxBackups int    `json:"maxbackups" bson:"maxbackups" validate:"gte=0"`
	LocalTime  bool   `json:"localtime"  bson:"localtime"`
	Compress   bool   `json:"compress"   bson:"compress"`
}
