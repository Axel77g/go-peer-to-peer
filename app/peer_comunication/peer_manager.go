package peer_comunication

import (
	"log"
	"sync"
)

var (
    peers = make(map[string]IPeer) //ip string to IPeer mapping
    peersMutex = sync.RWMutex{}
)

func GetPeerByAddress(address TransportAddress)  IPeer {
	ip := address.GetIP().String()
	peersMutex.RLock()
    defer peersMutex.RUnlock()
	if peer, exists := peers[ip]; exists {
		return peer
	}
	return nil
}

func RegisterTransportChannel(channel ITransportChannel) IPeer {
	address := channel.GetAddress()
	ip := address.GetIP().String()
	log.Printf("Registering transport channel for address: %s\n", address.String())

	peersMutex.RLock()
    defer peersMutex.RUnlock()

	if peer, exists := peers[ip]; exists {
		peer.addTransportChannel(channel)
		return peer
	}
	newPeer := NewPeer(address.GetIP())
	newPeer.addTransportChannel(channel)
	AddPeer(newPeer)
	return newPeer
}

func UnregisterTransportChannel(channel ITransportChannel) {
	address := channel.GetAddress()
	ip := address.GetIP().String()
	log.Printf("Removing transport channel for address: %s\n", address.String())

	peersMutex.Lock()
	defer peersMutex.Unlock()

	if peer, exists := peers[ip]; exists {
		peer.removeTransportChannel(channel)
		if len(peer.getTransportsChannels()) == 0 {
			delete(peers, ip)
			log.Printf("Peer %s removed due to no active transport channels\n", ip)
		}
	} else {
		log.Printf("No peer found for address: %s\n", address.String())
	}
}

func AddPeer(peer IPeer) {
	ip := peer.getAddress().String()
	if _, exists := peers[ip]; !exists {
		peers[ip] = peer
	}
	//for each peer log 
	for _, peer := range peers {
		log.Printf("[PMANAGE] Peer : %s", peer.String())
	}
}