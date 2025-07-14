package peer_comunication

import (
	"log"
	"peer-to-peer/app/shared"
	"sync"
)

var (
	peers      = make(map[string]IPeer) //ip string to IPeer mapping
	peersMutex = sync.Mutex{}
)

func GetPeerByAddress(address TransportAddress) (IPeer, bool) {
	ip := address.GetIP().String()
	peersMutex.Lock()
	defer peersMutex.Unlock()
	if peer, exists := peers[ip]; exists {
		return peer, true
	}
	return nil, false
}

func RegisterTransportChannel(channel ITransportChannel) IPeer {
	address := channel.GetAddress()
	ip := address.GetIP().String()
	log.Printf("Register %s transport channel for address: %s\n", channel.GetProtocol(), address.String())

	peersMutex.Lock()
	defer peersMutex.Unlock()

	if peer, exists := peers[ip]; exists {
		peer.AddTransportChannel(channel)
		debug()
		return peer
	}
	newPeer := NewPeer(address.GetIP())
	newPeer.AddTransportChannel(channel)
	AddPeer(newPeer)
	debug()
	return newPeer
}

func UnregisterTransportChannel(channel ITransportChannel) {
	address := channel.GetAddress()
	ip := address.GetIP().String()
	log.Printf("Unregister transport channel for address: %s\n", address.String())

	peersMutex.Lock()
	defer peersMutex.Unlock()

	if peer, exists := peers[ip]; exists {
		peer.RemoveTransportChannel(channel)
		if len(peer.GetTransportsChannels().channels) == 0 {
			delete(peers, ip)
			log.Printf("Peer %s removed due to no active transport channels\n", ip)
		}
	} else {
		log.Printf("No peer found for address: %s\n", address.String())
	}

	debug()
}

func BroadcastMessage(message []byte, protocol string) {
	peersMutex.Lock()
	defer peersMutex.Unlock()

	for _, peer := range peers {
		channels := peer.GetTransportsChannels()
		channel, exists := channels.GetByType(protocol)
		if exists {
			channel.Send(message)
		}
	}
}

func BroadcastIterator(message []byte, iterator shared.Iterator, protocol string) {
	peersMutex.Lock()
	defer peersMutex.Unlock()

	for _, peer := range peers {
		channels := peer.GetTransportsChannels()
		channel, exists := channels.GetByType(protocol)
		if exists {
			channel.SendIterator(message, iterator)
		}
	}
}

func AddPeer(peer IPeer) {
	ip := peer.GetIP().String()
	if _, exists := peers[ip]; !exists {
		peers[ip] = peer
	}
}

func debug() {
	//for each peer log
	for _, peer := range peers {
		log.Printf("[PMANAGE] Peer : %s", peer.String())
	}
}
