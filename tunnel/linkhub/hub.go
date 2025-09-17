package linkhub

import (
	"sync"

	"github.com/xmx/aegis-common/transport"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Huber[K comparable] interface {
	// Get 通过 ID 获取 peer。
	Get(id K) transport.Peer[K]

	// Del 通过删除节点，并返回删除的数据，nil 代表之前无数据。
	Del(id K) transport.Peer[K]

	// Put 存放节点并返回原来的数据（如果存在的话）。
	Put(p transport.Peer[K]) (old transport.Peer[K])

	// PutIfAbsent 当此 ID 没有时才存放节点，并返回是否放入成功。
	PutIfAbsent(p transport.Peer[K]) bool
}

func NewHub(capacity int) Huber[bson.ObjectID] {
	if capacity < 0 {
		capacity = 0
	}

	return &simpleHub[bson.ObjectID]{
		peers: make(map[bson.ObjectID]transport.Peer[bson.ObjectID], capacity),
	}
}

type simpleHub[K comparable] struct {
	mutex sync.RWMutex
	peers map[K]transport.Peer[K]
}

func (sh *simpleHub[K]) Get(id K) transport.Peer[K] {
	sh.mutex.RLock()
	defer sh.mutex.RUnlock()

	return sh.peers[id]
}

func (sh *simpleHub[K]) Del(id K) transport.Peer[K] {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	old := sh.peers[id]
	delete(sh.peers, id)

	return old
}

func (sh *simpleHub[K]) Put(p transport.Peer[K]) transport.Peer[K] {
	id := p.ID()

	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	old := sh.peers[id]
	sh.peers[id] = p

	return old
}

func (sh *simpleHub[K]) PutIfAbsent(p transport.Peer[K]) bool {
	id := p.ID()

	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	_, exists := sh.peers[id]
	if !exists {
		sh.peers[id] = p
	}

	return !exists
}
