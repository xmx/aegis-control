package linkhub

import (
	"github.com/xmx/aegis-common/tunnel/tundial"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Peer interface {
	// ID 节点数据库 ID。
	ID() bson.ObjectID

	// Muxer 底层通道。
	Muxer() tundial.Muxer
}

type Peers []Peer

func NewPeer(id bson.ObjectID, mux tundial.Muxer) Peer {
	return &muxerPeer{
		id:  id,
		mux: mux,
	}
}

type muxerPeer struct {
	id  bson.ObjectID
	mux tundial.Muxer
}

func (p *muxerPeer) ID() bson.ObjectID {
	return p.id
}

func (p *muxerPeer) Muxer() tundial.Muxer {
	return p.mux
}
