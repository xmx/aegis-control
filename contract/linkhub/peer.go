package linkhub

import (
	"github.com/xmx/aegis-common/transport"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Peer interface {
	Host() string
	Muxer() transport.Muxer
	ObjectID() bson.ObjectID
}

func NewPeer(oid bson.ObjectID, mux transport.Muxer) Peer {
	return &tunnelPeer{
		oid:  oid,
		mux:  mux,
		host: oid.Hex(),
	}
}

type tunnelPeer struct {
	oid  bson.ObjectID
	mux  transport.Muxer
	host string
}

func (t *tunnelPeer) Host() string {
	return t.host
}

func (t *tunnelPeer) Muxer() transport.Muxer {
	return t.mux
}

func (t *tunnelPeer) ObjectID() bson.ObjectID {
	return t.oid
}
