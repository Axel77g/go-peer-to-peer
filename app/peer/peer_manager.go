package peer

import (
	"fmt"
	"sync"
)

type Peer struct {
	ID   string
	Addr string
}

type PeerManager struct {
	peers   map[string]Peer
	updates chan Peer
	mu      sync.Mutex
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers:   make(map[string]Peer),
		updates: make(chan Peer),
	}
}

func (pm *PeerManager) AddPeer(peer Peer) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.peers[peer.ID] = peer
	fmt.Printf("Peer ajouté/actualisé : %v\n", peer)
}

func (pm *PeerManager) Has(peer Peer) bool {
	//check if the map as the peer ID
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, ok := pm.peers[peer.ID]; ok {
		return true
	} else {
		return false
	}
}
