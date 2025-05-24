package main

import (
	"log"
	"net"
	"peer-to-peer/app/discovery"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
	"strconv"
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
	go func(){ // Handle new udp transport channel
		select {
			case channel := <- udpServer.NewTransportChannelEvent:
				peer_comunication.RegisterTransportChannel(channel) //register the new udp transport channel
				//try connect to the peer using the tcp protocol
				addr := channel.GetAddress()
				conn, err := net.Dial("tcp", net.JoinHostPort(addr.GetIP().String(), strconv.Itoa(shared.TCPPort)))
				if( err != nil) {
					log.Println("Error connecting to peer via TCP:", err)
					return
				}
				transportChannel := peer_comunication.NewTCPTransportChannel(conn)
				peer_comunication.RegisterTransportChannel(transportChannel) //register the new local tcp transport channel
		}
	}()

	tcpServer := peer_comunication.NewTCPServer(shared.TCPPort)
	go tcpServer.Listen() //open tcp server listener
	go func() { // Handle new tcp transport channel
		select {
			case channel := <- tcpServer.NewTransportChannelEvent:
				peer_comunication.RegisterTransportChannel(channel) //register the new remote tcp transport channel
		}
	}()

	networkInterfaceManager := discovery.NewNetworkInterfaceManager()
	go discovery.SenderLoop(shared.SOCKET_ID, networkInterfaceManager)

    wg.Wait()
}