package peer

import (
	"fmt"
	"sync"
)

type Peer struct {
	ID   string
	Addr string
}

func ManagePeers(peerUpdates <-chan Peer) {
	peers := make(map[string]Peer)
	var mu sync.Mutex

	for peer := range peerUpdates {
		mu.Lock()
		peers[peer.ID] = peer
		mu.Unlock()
		fmt.Printf("Peer ajouté/actualisé : %v\n", peer)
	}
}
