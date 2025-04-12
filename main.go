package main

import (
	"log"
	"math/rand"
	"peer-to-peer/app/discovery"
	filesystemwatcher "peer-to-peer/app/file_system_watcher"
	filetransfering "peer-to-peer/app/file_transfering"
	"peer-to-peer/app/peer"
	"sync"
	"time"
)

func main() {
	const SHARED_DIRECTORY = "./Shared"

	wg := sync.WaitGroup{}
	wg.Add(2)

	SOCKET_ID := rand.Intn(10000000)

	peerManager := peer.NewPeerManager()
	transferingQueue := filetransfering.NewTransferQueue()
	networkInterfaceManager := discovery.NewNetworkInterfaceManager()

	fileEvents := make(chan filesystemwatcher.FileSystemEvent)
	fileWatcher := filesystemwatcher.NewWatcher(SHARED_DIRECTORY, 2*time.Second, fileEvents)
	go fileWatcher.Listen()
	go func() {
		for event := range fileEvents {
			log.Println("Event reçu:", event.EventType, "sur le fichier:", event.FilePath)
		}
	}()
	go discovery.ListenForDiscoverRequests(SOCKET_ID, peerManager)
	go discovery.SenderLoop(SOCKET_ID, networkInterfaceManager, peerManager)
	go transferingQueue.Loop()

	wg.Wait()
}
