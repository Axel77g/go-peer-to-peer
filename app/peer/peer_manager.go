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
	pm.peers[peer.ID] = peer
	pm.mu.Unlock()
	pm.PrintPeer()
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
		log.Printf("Peer : %s, Addr : %s\n", peer.ID, peer.Addr.String())
	}
}