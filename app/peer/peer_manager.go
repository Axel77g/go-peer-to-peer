package peer

import (
	"log"
	"sync"
	"time"
)

type PeerManager struct {
	peers   map[string]Peer
	Updates chan Peer
	mu      sync.Mutex
}

var instance *PeerManager
var once sync.Once

func GetPeerManager() *PeerManager {
	once.Do(func() {
		instance = NewPeerManager()
	})
	return instance
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers:   make(map[string]Peer),
		Updates: make(chan Peer),
	}
}

func (pm *PeerManager) UpsertPeer(peer Peer) {
	pm.mu.Lock()
	_, exists := pm.peers[peer.ID]
	pm.peers[peer.ID] = peer
	pm.mu.Unlock()
	if !exists { pm.Updates <- peer }
}

func (pm *PeerManager) RemoveInactivePeers(timeout time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for id, peer := range pm.peers {
		if time.Since(peer.LastSeen) > timeout {
			delete(pm.peers, id)
			log.Printf("Peer %s removed due to inactivity\n", id)
		}
	}
}

func (pm *PeerManager) PrintPeer() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for _, peer := range pm.peers {
		log.Printf("[PMANAGE] Peer : %s, Addr : %s -> Socket  %v \n ", peer.ID, peer.Addr.String(), peer.TCPSocket)
	}
}

func (pm *PeerManager) GetPeer(id string) (Peer, bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	peer, exists := pm.peers[id]
	return peer, exists
}