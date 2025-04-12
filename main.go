package main

import (
	"math/rand"
	"peer-to-peer/app/discovery"
	filesystemwatcher "peer-to-peer/app/file_system_watcher"
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

	fileEvents := make(chan filesystemwatcher.FileSystemEvent)
	fileWatcher := filesystemwatcher.NewWatcher(SHARED_DIRECTORY, 2*time.Second, fileEvents)
	go fileWatcher.Listen()
	go func() {
		for event := range fileEvents {
			println("Event reçu:", event.EventType, "sur le fichier:", event.FilePath)
		}
	}()
	go discovery.ListenForDiscoverRequests(SOCKET_ID, peerManager)
	go discovery.SenderLoop(SOCKET_ID)
	go transferingQueue.Loop()

	wg.Wait()
}
