package discovery

import (
	"fmt"
	"log"
	"net"
	"peer-to-peer/app/peer"
	"time"
)

/*
Emit on the broadcast network address
*/
func discoveryRequestSender(socketID int, ip net.IP) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ip,
		Port: 9999,
	})

	if err != nil {
		log.Printf("Erreur de connexion UDP broadcast : %v\n", err)
		return
	}
	defer conn.Close()

	// Autoriser le broadcast (pas obligatoire dans Go, mais bonne pratique)
	if err := conn.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
		log.Println("Erreur deadline UDP :", err)
		return
	}

	message := []byte("DISCOVER_PEER_REQUEST:" + fmt.Sprintf("%d", socketID))
	_, err = conn.Write(message)
	if err != nil {
		log.Printf("Erreur lors de l'envoi UDP : %v\n", err)
		return
	}
}

func SenderLoop(socketID int, networkInterfaceManager *NetworkInterfaceManager, peerManager *peer.PeerManager) {
	networkInterfaceManager.fetchInterfaces()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, ip := range networkInterfaceManager.AvailableIpInterface {
				broadcastAddr := ip.getBroadcastAddress()
				discoveryRequestSender(socketID, broadcastAddr)
				networkInterfaceManager.fetchInterfaces()
			}
			peerManager.RemoveInactivePeers(10 * time.Second)
		}
	}
}
