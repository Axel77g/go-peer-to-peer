package peer

import (
	"net"
	"peer-to-peer/app/shared"
	"time"
)

type Peer struct {
	ID       string
	Addr     net.IP
	TCPSocket shared.Socket
	LastSeen time.Time
}

func NewPeer(id string, addr net.IP) Peer {
	return Peer{
		ID:       id,
		Addr:     addr,
		LastSeen: time.Now(),
	}
}