package main

import (
	"log"
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

	
	log.Println("Socket ID:", shared.SOCKET_ID)

	tcpServer := tcpcomunication.NewTCPServer()
	peerManager := peer.GetPeerManager()
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
			log.Println("Nouvelle peer détecté:", peer.ID)
			if peer.TCPSocket != nil {
				log.Println("Cette nouvelle paire a déjà une connexion TCP active:", peer.ID)
				continue
			}
			socket,err := tcpcomunication.CreateTCPConnection(&peer)
			if err != nil {
				log.Println("Erreur de connexion TCP au peer:", err)
				continue
			}
			peer.TCPSocket = &socket
			peerManager.UpdatePeer(peer)
		}
	}()

	go tcpServer.Listen()
	go fileWatcher.Listen()
	go discovery.Listen(shared.SOCKET_ID, peerManager)
	
	go discovery.SenderLoop(shared.SOCKET_ID, networkInterfaceManager, peerManager)
	go transferingQueue.Loop()

	wg.Wait()
}
