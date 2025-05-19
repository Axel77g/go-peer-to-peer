package peer_comunication

import "net"

type IPeer interface {
	getHash() string
	getAddress() net.Addr
	getTransportsChannel() []ITransportChannel
}