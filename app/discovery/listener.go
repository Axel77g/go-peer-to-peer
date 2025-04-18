package discovery

import (
	"fmt"
	"log"
	"net"
	"peer-to-peer/app/peer"
	"peer-to-peer/app/shared"
	"strings"
)

func Listen(socketID int, pm *peer.PeerManager) {
	addr := net.UDPAddr{
		Port: shared.UDPPort,
		IP:   net.IPv4zero,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Printf("Erreur lors de l'écoute UDP : %v\n", err)
		return
	}
	defer conn.Close()

	log.Println("Serveur de découverte UDP en écoute sur le port 9999")

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Erreur de lecture UDP : %v\n", err)
			continue
		}

		message := strings.TrimSpace(string(buffer[:n]))
		//split le message en ":"
		parts := strings.Split(message, ":")
		if len(parts) > 1 {
			if parts[1] == fmt.Sprintf("%d", socketID) {
				continue
			}
			if parts[0] == "DISCOVER_PEER_REQUEST" {
				peerInstance, exist := pm.GetPeer(parts[1])
				if !exist {
					peerInstance = peer.NewPeer(
						parts[1], 
						remoteAddr.IP,
					)
				}else{
					peerInstance.Signal()
				}
				pm.UpsertPeer(peerInstance)
			}
		}

	}
}
