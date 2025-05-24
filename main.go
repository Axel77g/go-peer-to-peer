package main

import (
	"peer-to-peer/app/discovery"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
	"sync"
)

/* func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)


	log.Println("Socket ID:", shared.SOCKET_ID)

	tcpServer := tcpcomunication.NewTCPServer()
	peerManager := peer.GetPeerManager()
	transferingQueue := filetransfering.NewTransferQueue()
	networkInterfaceManager := discovery.NewNetworkInterfaceManager()
	fileEvents := make(chan file_event.FileEvent)
	fileWatcher := file_watcher.NewWatcher(shared.SHARED_DIRECTORY, 2*time.Second, fileEvents)

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
			peerManager.UpsertPeer(peer)
		}
	}()

	go tcpServer.Listen()
	go fileWatcher.Listen()
	go discovery.Listen(shared.SOCKET_ID, peerManager)

	go discovery.SenderLoop(shared.SOCKET_ID, networkInterfaceManager, peerManager)
	go transferingQueue.Loop()

	wg.Wait()
}
*/


func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	udpServer := peer_comunication.GetUDPServerListener()
	go udpServer.Listen() //open udp server listener
	tcpServer := peer_comunication.NewTCPServer(shared.TCPPort)
	go tcpServer.Listen() //open tcp server listener

	networkInterfaceManager := discovery.NewNetworkInterfaceManager()
	go discovery.SenderLoop(shared.SOCKET_ID, networkInterfaceManager)

    wg.Wait()
}