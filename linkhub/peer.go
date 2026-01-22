package linkhub

import (
	"context"
	"time"

	"github.com/xmx/aegis-common/muxlink/muxconn"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ConfigLoader[T any] interface {
	LoadConfig(ctx context.Context) (*T, error)
}

type ServerHooker interface {
	OnConnected(inf Info, connectAt time.Time)

	OnDisconnected(inf Info, connectAt, disconnectAt time.Time)
}

type Peer interface {
	// ID 节点数据库 ID。
	ID() bson.ObjectID

	// Host 主机名。
	Host() string

	// Muxer 底层通道。
	Muxer() muxconn.Muxer

	// Info 节点信息。
	Info() Info
}

type muxPeer struct {
	id   bson.ObjectID
	mux  muxconn.Muxer
	inf  Info
	host string
}

func (m *muxPeer) ID() bson.ObjectID    { return m.id }
func (m *muxPeer) Host() string         { return m.host }
func (m *muxPeer) Muxer() muxconn.Muxer { return m.mux }
func (m *muxPeer) Info() Info           { return m.inf }

type Info struct {
	Name     string `json:"name"`
	Inet     string `json:"inet"`
	Goos     string `json:"goos"`
	Goarch   string `json:"goarch"`
	Hostname string `json:"hostname"`
	Semver   string `json:"semver"`
}
