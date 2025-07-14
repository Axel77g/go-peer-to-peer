package main

import (
	"peer-to-peer/app/discovery"
	file_event "peer-to-peer/app/files/event"
	file_watcher "peer-to-peer/app/files/watcher"
	"peer-to-peer/app/handlers"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	tcpServer := peer_comunication.NewTCPServer(shared.TCPPort)
	go tcpServer.Listen(&handlers.TCPControllerTransportChannelHandler{}) //open tcp server listener with the tcp handler

	udpServer := peer_comunication.GetUDPServerListener()
	go udpServer.Listen(&handlers.UDPDiscoveryTransportChannel{}) //open udp server listener with the discovery handler

	networkInterfaceManager := discovery.NewNetworkInterfaceManager()
	go discovery.SenderLoop(shared.SOCKET_ID, networkInterfaceManager)

	eventManager := file_event.GetEventManager()

	dirChan := make(chan shared.FileEvent, 100)
	watcher := file_watcher.NewWatcher(shared.SHARED_DIRECTORY, time.Second*5, dirChan)
	go watcher.Listen() //start the file watcher

	go func() {
		for event := range dirChan {
			eventManager.AppendEvent(event)
			eventManager.BroadcastEvents()
		}
	}()

	wg.Wait()
}
