package peer_comunication

import "net"

type TransportAddress struct {
	ip   net.IP
	port int
}
