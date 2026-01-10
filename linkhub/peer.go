package linkhub

import (
	"time"

	"github.com/xmx/aegis-common/muxlink/muxconn"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Peer interface {
	// ID 节点数据库 ID。
	ID() bson.ObjectID

	// Muxer 底层通道。
	Muxer() muxconn.Muxer

	// Info is node info.
	Info() Info
}

type Peers []Peer

func NewPeer(id bson.ObjectID, mux muxconn.Muxer, info Info) Peer {
	return &muxerPeer{
		id:   id,
		mux:  mux,
		info: info,
	}
}

type muxerPeer struct {
	id   bson.ObjectID
	mux  muxconn.Muxer
	info Info
}

func (p *muxerPeer) ID() bson.ObjectID    { return p.id }
func (p *muxerPeer) Muxer() muxconn.Muxer { return p.mux }
func (p *muxerPeer) Info() Info           { return p.info }

type Info struct {
	Name     string `json:"name"`
	Inet     string `json:"inet"`
	Goos     string `json:"goos"`
	Goarch   string `json:"goarch"`
	Hostname string `json:"hostname"`
	Semver   string `json:"semver"`
}

// ConnectListener 节点上下线通知接口。
type ConnectListener interface {
	// OnConnection 节点上线事件通知。
	//
	// 该方法为同步阻塞方法。
	OnConnection(peer Peer, connectAt time.Time)

	// OnDisconnection 节点下线事件通知。
	//
	// 该方法为同步阻塞方法。
	OnDisconnection(info Info, connectAt, disconnectAt time.Time)
}
