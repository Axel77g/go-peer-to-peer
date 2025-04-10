package discovery

import (
	"fmt"
	"net"
	"peer-to-peer/app/peer"
	"strings"
)

func ListenForDiscoverRequests(socketID int, peerUpdates chan<- peer.Peer) {
	addr := net.UDPAddr{
		Port: 9999,
		IP:   net.ParseIP("0.0.0.0"), // écoute sur toutes les interfaces
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Erreur lors de l'écoute UDP : %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Écoute UDP sur le port 9999...")

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Erreur de lecture UDP : %v\n", err)
			continue
		}

		message := strings.TrimSpace(string(buffer[:n]))
		//split le message en ":"
		parts := strings.Split(message, ":")
		if len(parts) > 1 {
			if parts[1] == fmt.Sprintf("%d", socketID) {
				fmt.Println("-> Ignore le message")
				continue
			}
			if parts[0] == "DISCOVER_PEER_REQUEST" {
				fmt.Printf("Message DISCOVER_PEER_REQUEST reçu de %s\n", remoteAddr.String())
				peer := peer.Peer{
					ID:   parts[1],
					Addr: remoteAddr.String(),
				}
				peerUpdates <- peer                             // Send peer data to the channel
				discoveryRequestSender(socketID, remoteAddr.IP) //send the discovery request to the peer
			}
		}

	}
}
