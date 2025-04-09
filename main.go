package main

import (
	"math/rand"
	"peer-to-peer/app/discovery"
	"peer-to-peer/app/peer"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	SOCKET_ID := rand.Intn(10000000)

	peerUpdates := make(chan peer.Peer)

	go peer.ManagePeers(peerUpdates)
	go discovery.ListenForDiscoverRequests(SOCKET_ID, peerUpdates)
	go discovery.SenderLoop(SOCKET_ID)
	wg.Wait()
}
