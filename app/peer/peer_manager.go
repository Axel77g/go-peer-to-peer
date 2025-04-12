package peer

import (
	"log"
	"sync"
	"time"
)

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

func (pm *PeerManager) SignalPeer(peer Peer) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.peers[peer.ID] = peer
	log.Printf("Peer ajouté/Actualisé : %v\n", peer)
}

func (pm *PeerManager) RemoveInactivePeers(timeout time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for id, peer := range pm.peers {
		if time.Since(peer.LastSeen) > timeout {
			delete(pm.peers, id)
			log.Printf("Peer supprimé : %v\n", peer)
		}
	}
}
