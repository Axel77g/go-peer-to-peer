package peer

import (
	"net"
	"time"
)

type Peer struct {
	ID       string
	Addr     net.IP
	TCPPort  uint16
	LastSeen time.Time
}

func NewPeer(id string, addr net.IP) Peer {
	return Peer{
		ID:       id,
		Addr:     addr,
		LastSeen: time.Now(),
	}
}