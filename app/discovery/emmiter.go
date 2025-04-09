package discovery

import (
	"fmt"
	"net"
	"time"
)

/*
Emit on the broadcast network address
*/
func discoveryRequestSender(socketID int) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4bcast, // 255.255.255.255
		Port: 9999,
	})

	if err != nil {
		fmt.Printf("Erreur de connexion UDP broadcast : %v\n", err)
		return
	}
	defer conn.Close()

	// Autoriser le broadcast (pas obligatoire dans Go, mais bonne pratique)
	if err := conn.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
		fmt.Println("Erreur deadline UDP :", err)
		return
	}

	message := []byte("DISCOVER_PEER_REQUEST:" + fmt.Sprintf("%d", socketID))
	_, err = conn.Write(message)
	if err != nil {
		fmt.Printf("Erreur lors de l'envoi UDP : %v\n", err)
		return
	}

	fmt.Println("→ Message DISCOVER_PEER_REQUEST émis en broadcast")
}

func SenderLoop(socketID int) {
	discoveryRequestSender(socketID)

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Erreur lors de la récupération des interfaces réseau : %v\n", err)
		return
	}

	previousAddrs := make(map[string][]string)

	for {
		for _, iface := range interfaces {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Printf("Erreur lors de la récupération des adresses pour l'interface %s : %v\n", iface.Name, err)
				continue
			}

			currentAddrs := []string{}
			for _, addr := range addrs {
				currentAddrs = append(currentAddrs, addr.String())
			}

			// Check for changes in IP addresses
			if prev, exists := previousAddrs[iface.Name]; !exists || !equalSlices(prev, currentAddrs) {
				fmt.Printf("Changement détecté sur l'interface %s\n", iface.Name)
				previousAddrs[iface.Name] = currentAddrs
				discoveryRequestSender(socketID)
			}
		}

		time.Sleep(5 * time.Second) // Adjust the interval as needed
	}
}

// Helper function to compare slices
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
