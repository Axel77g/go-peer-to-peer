package discovery

import (
	"fmt"
	"log"
	"net"
	"peer-to-peer/app/handlers"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
	"time"
)

/*
Emit on the broadcast network address
*/
func discoveryRequestSender(socketID int, ip net.IP) {
	address := peer_comunication.NewTransportAddress(ip, shared.UDPPort)
	transportChannel := peer_comunication.NewUDPTransportChannel(address, &handlers.UDPDiscoveryEmitterTransportChannelHandler{})
	if transportChannel == nil {
		log.Println("Failed to create UDP transport channel")
		return
	}
		
	message := []byte("DISCOVER_PEER_REQUEST:" + fmt.Sprintf("%d", socketID))
	err := transportChannel.Send(message)
	if err != nil {
		log.Printf("Failed to send discovery request to %s: %v\n", address.String(), err)
		return
	}
}

func sender(networkInterfaceManager *NetworkInterfaceManager, socketID int) {
	log.Printf("Sending discovery request to all available IP interfaces with socket ID %d\n", socketID)
	for _, ip := range networkInterfaceManager.AvailableIpInterface {
		broadcastAddr := ip.getBroadcastAddress()
		discoveryRequestSender(socketID, broadcastAddr)
	}
}

func SenderLoop(socketID int, networkInterfaceManager *NetworkInterfaceManager) {
	networkInterfaceManager.fetchInterfaces()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	sender(networkInterfaceManager, socketID)
	for range ticker.C {
		sender(networkInterfaceManager, socketID)
	}
}
