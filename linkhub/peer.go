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

func NewPeer(id bson.ObjectID, mux tundial.Muxer) Peer {
	return &peer{
		id:  id,
		mux: mux,
	}
}

type peer struct {
	id  bson.ObjectID
	mux tundial.Muxer
}

func (p peer) ID() bson.ObjectID {
	return p.id
}

func (p peer) Muxer() tundial.Muxer {
	return p.mux
}
