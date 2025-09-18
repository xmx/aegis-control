package linkhub

import (
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrDuplicateConnection = errors.New("节点重复连接上线")

type Huber interface {
	// Get 获取节点。
	Get(host string) Peer
	GetByObjectID(id bson.ObjectID) Peer

	// Del 删除节点，并返回删除的数据，nil 代表之前无数据。
	Del(host string) Peer

	// Put 不存在时才存放节点，并返回是否放入成功。
	Put(p Peer) bool
}

func NewHub(capacity int) Huber {
	if capacity < 0 {
		capacity = 0
	}

	return &hashmapHub{
		peers: make(map[string]Peer, capacity),
	}
}

type hashmapHub struct {
	mutex sync.RWMutex
	peers map[string]Peer
}

func (h *hashmapHub) Get(host string) Peer {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.peers[host]
}

func (h *hashmapHub) GetByObjectID(id bson.ObjectID) Peer {
	return h.Get(id.Hex())
}

func (h *hashmapHub) Del(host string) Peer {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	old := h.peers[host]
	delete(h.peers, host)

	return old
}

func (h *hashmapHub) Put(p Peer) bool {
	host := p.Host()
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, exists := h.peers[host]
	if !exists {
		h.peers[host] = p
	}

	return !exists
}
