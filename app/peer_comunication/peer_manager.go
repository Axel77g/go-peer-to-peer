package peer_comunication

import (
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

func AddPeer(peer IPeer) {
	ip := peer.getAddress().String()
	if _, exists := peers[ip]; !exists {
		peers[ip] = peer
	}
}