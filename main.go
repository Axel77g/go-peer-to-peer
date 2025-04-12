package main

import (
	"log"
	"math/rand"
	"peer-to-peer/app/discovery"
	filesystemwatcher "peer-to-peer/app/file_system_watcher"
	filetransfering "peer-to-peer/app/file_transfering"
	"peer-to-peer/app/peer"
	"peer-to-peer/app/shared"
	tcpcomunication "peer-to-peer/app/tcp_comunication"
	"sync"
	"time"
)


func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	SOCKET_ID := rand.Intn(10000000)

	tcpServer := tcpcomunication.NewTCPServer()
	peerManager := peer.NewPeerManager()
	transferingQueue := filetransfering.NewTransferQueue()
	networkInterfaceManager := discovery.NewNetworkInterfaceManager()
	fileEvents := make(chan filesystemwatcher.FileSystemEvent)
	fileWatcher := filesystemwatcher.NewWatcher(shared.SHARED_DIRECTORY, 2*time.Second, fileEvents)
	
	go func() {
		for event := range fileEvents {
			log.Println("Event reçu:", event.EventType, "sur le fichier:", event.FilePath)
		}
	}()

	go func(){
		for peer := range peerManager.Updates {
			//connect to the peer in TCP
			tcpClient := tcpcomunication.NewTCPClient(peer)
			if err := tcpClient.Connect(); err != nil {
				log.Println("Erreur de connexion au peer:", err)
				continue
			}
		}
	}()

	go tcpServer.Listen()
	go fileWatcher.Listen()
	go discovery.Listen(SOCKET_ID, peerManager)
	
	go discovery.SenderLoop(SOCKET_ID, networkInterfaceManager, peerManager)
	go transferingQueue.Loop()

	wg.Wait()
}
