package linkhub

import (
	"sync"

	"github.com/xmx/aegis-common/muxlink/muxconn"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Huber interface {
	// Put 将节点加入到连接池，如果 id 已存在则加入失败，返回 nil。
	// 加入成功则返回 Peer 节点。
	Put(id bson.ObjectID, mux muxconn.Muxer, inf Info) Peer

	Get(host string) Peer

	GetID(id bson.ObjectID) Peer

	Del(host string) Peer

	DelID(id bson.ObjectID) Peer

	// Domain 域。
	Domain() string

	Peers() []Peer
}

func NewHub(domain string) Huber {
	return &safeMapHub{
		domain: domain,
		peers:  make(map[string]Peer, 16),
	}

}

type safeMapHub struct {
	domain string
	mutex  sync.RWMutex
	peers  map[string]Peer
}

func (s *safeMapHub) Put(id bson.ObjectID, mux muxconn.Muxer, inf Info) Peer {
	host := resolveHost(id, s.domain)
	peer := &muxPeer{
		id:   id,
		mux:  mux,
		inf:  inf,
		host: host,
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.peers[host]; exists {
		return nil
	}
	s.peers[host] = peer

	return peer
}

func (s *safeMapHub) Get(host string) Peer {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.peers[host]
}

func (s *safeMapHub) GetID(id bson.ObjectID) Peer {
	host := resolveHost(id, s.domain)
	return s.Get(host)
}

func (s *safeMapHub) Del(host string) Peer {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	peer := s.peers[host]
	if peer != nil {
		delete(s.peers, host)
	}

	return peer
}

func (s *safeMapHub) DelID(id bson.ObjectID) Peer {
	host := resolveHost(id, s.domain)
	return s.Del(host)
}

func (s *safeMapHub) Domain() string {
	return s.domain
}

func (s *safeMapHub) Peers() []Peer {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	res := make([]Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		res = append(res, peer)
	}

	return res
}

func resolveHost(id bson.ObjectID, domain string) string {
	host := id.Hex()
	return host + "." + domain
}
