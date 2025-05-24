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

/**
 * this is custom business logic to handle the new remote transport channel
 * register a channel will create a peer in the peer manager if it does not exist and bind the channel to the peer
 * then it will try to connect to the peer using the tcp protocol to establish a tcp connection
 */
func handleNewRemoteUDPTransportChannel(udpServer *peer_comunication.UDPServerListener) {
	for {
		select {
			case channel := <- udpServer.NewTransportChannelEvent:
				peer_comunication.RegisterTransportChannel(channel) //register the new udp transport channel -> this will create a peer if it does not exist and trigger handleNewPeerAddedEvent
		}
	}
}

/**
 * this is custom business logic to handle the new remote transport channel from TCP server
 * register a channel will create a peer in the peer manager if it does not exist and bind the channel to the peer
 */
func hanldeNewRemoteTCPTransportChannel(tcpServer *peer_comunication.TCPServer) {
	for {
		select {
			case channel := <- tcpServer.NewTransportChannelEvent:
				peer_comunication.RegisterTransportChannel(channel) //register the new remote tcp transport channel
		}
	}
}

func handleNewPeerAddedEvent(){
	for {
		select {
		case peer := <- peer_comunication.NewPeerAddedUpdate:
			channelCollection := peer.GetTransportsChannels()
			_, exist := channelCollection.GetByType("tcp")
			if !exist { //if the tcp not exist, we will try to connect to the peer using the tcp protocol
				conn, err := net.Dial("tcp", net.JoinHostPort(peer.GetIP().String(), strconv.Itoa(shared.TCPPort)))
				if( err != nil) {
					log.Println("Error connecting to peer via TCP:", err)
					return
				}
				transportChannel := peer_comunication.NewTCPTransportChannel(conn)
				peer_comunication.RegisterTransportChannel(transportChannel)
			}
		}
	}
}


func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	udpServer := peer_comunication.GetUDPServerListener()
	go udpServer.Listen() //open udp server listener

	tcpServer := peer_comunication.NewTCPServer(shared.TCPPort)
	go tcpServer.Listen() //open tcp server listener

	go handleNewRemoteUDPTransportChannel(udpServer) // Handle new udp transport channel
	go handleNewPeerAddedEvent()
	go hanldeNewRemoteTCPTransportChannel(tcpServer) // Handle new tcp transport channel

	networkInterfaceManager := discovery.NewNetworkInterfaceManager()
	go discovery.SenderLoop(shared.SOCKET_ID, networkInterfaceManager)

    wg.Wait()
}