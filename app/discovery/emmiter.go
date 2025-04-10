package discovery

import (
	"fmt"
	"net"
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
}

func SenderLoop(socketID int) {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.To4() == nil {
				continue // ignore IPv6 ou adresses non-IPNet
			}

			ip := ipnet.IP.To4()
			mask := ipnet.Mask

			broadcast := make(net.IP, 4)
			for i := 0; i < 4; i++ {
				broadcast[i] = ip[i] | ^mask[i]
			}

			fmt.Printf("IP locale : %s | IP broadcast %s\n", ip.String(), broadcast.String())
			discoveryRequestSender(socketID, broadcast)
		}
	}
}
