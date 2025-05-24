package peer_comunication

import (
	"net"
	"strconv"
)

type TransportAddress struct {
	ip   net.IP
	port int
}
type TransportAddressKey struct {
	ip   string
	port int
}

func NewTransportAddress(ip net.IP, port int) TransportAddress {
	return TransportAddress{
		ip:   ip,
		port: port,
	}
}

func (t *TransportAddress) GetKey() TransportAddressKey {
	return TransportAddressKey{
		ip:   t.ip.String(),
		port: t.port,
	}
}

func (t *TransportAddress) GetIP() net.IP {
	return t.ip
}
func (t *TransportAddress) GetPort() int {
	return t.port
}

func (t *TransportAddress) String() string {
	return t.ip.String() + ":" + strconv.Itoa(t.port)
}