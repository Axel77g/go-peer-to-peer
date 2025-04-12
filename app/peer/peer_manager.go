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

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers:   make(map[string]Peer),
		Updates: make(chan Peer),
	}
}

func (pm *PeerManager) SignalPeer(peer Peer) {
	pm.mu.Lock()
	_, exists := pm.peers[peer.ID]
	pm.peers[peer.ID] = peer
	pm.mu.Unlock()
	if !exists { pm.Updates <- peer }
	//pm.PrintPeer()
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