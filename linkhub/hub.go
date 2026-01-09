package linkhub

import (
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Huber interface {
	Get(host string) Peer

	GetByID(id bson.ObjectID) Peer

	Put(p Peer) (succeed bool)

	Del(host string) Peer

	DelByID(id bson.ObjectID) Peer

	Peers() Peers
}

func NewHub(initSize ...int) Huber {
	size := 64
	if len(initSize) > 0 && initSize[0] > 0 {
		size = initSize[0]
	}

	return &safeMap{
		peers: make(map[string]Peer, size),
	}
}

type safeMap struct {
	mutex sync.RWMutex
	peers map[string]Peer
}

func (sm *safeMap) Get(host string) Peer {
	sm.mutex.RLock()
	peer := sm.peers[host]
	sm.mutex.RUnlock()

	return peer
}

func (sm *safeMap) GetByID(id bson.ObjectID) Peer {
	host := sm.toHost(id)
	return sm.Get(host)
}

func (sm *safeMap) Put(p Peer) bool {
	host := sm.toHost(p.ID())
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	_, exists := sm.peers[host]
	if exists {
		return false
	}

	sm.peers[host] = p

	return true
}

func (sm *safeMap) Del(host string) Peer {
	sm.mutex.Lock()
	peer := sm.peers[host]
	delete(sm.peers, host)
	sm.mutex.Unlock()

	return peer
}

func (sm *safeMap) DelByID(id bson.ObjectID) Peer {
	host := sm.toHost(id)
	return sm.Del(host)
}

func (sm *safeMap) Peers() Peers {
	return sm.snapshot()
}

func (*safeMap) toHost(id bson.ObjectID) string {
	return id.Hex()
}

func (sm *safeMap) snapshot() []Peer {
	sm.mutex.RLock()
	peers := make([]Peer, 0, len(sm.peers))
	for _, peer := range sm.peers {
		peers = append(peers, peer)
	}
	sm.mutex.RUnlock()

	return peers
}

type Huber1 interface {
	Put(p Peer) (succeed bool)
	GetByID(id bson.ObjectID) Peer
	GetByHost(host string) Peer
	DelByID(id bson.ObjectID) Peer
	DelByHost(host string) Peer
	Peers() Peers
}
