package main

import (
	"peer-to-peer/app/discovery"
	"peer-to-peer/app/handlers"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
	"sync"
)


func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	udpServer := peer_comunication.GetUDPServerListener()
	go udpServer.Listen(&handlers.UDPDiscoveryTransportChannel{}) //open udp server listener with the discovery handler

	tcpServer := peer_comunication.NewTCPServer(shared.TCPPort)
	go tcpServer.Listen(&handlers.TCPControllerTransportChannelHandler{}) //open tcp server listener with the tcp handler
	networkInterfaceManager := discovery.NewNetworkInterfaceManager()
	go discovery.SenderLoop(shared.SOCKET_ID, networkInterfaceManager)

    wg.Wait()
}