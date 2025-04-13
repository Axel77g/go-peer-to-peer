package peer

import (
	"net"
	"time"
)

type Peer struct {
	ID       string
	Addr     net.IP
	TCPSocket PeerMessager
	LastSeen time.Time
}

func NewPeer(id string, addr net.IP) Peer {
	return Peer{
		ID:       id,
		Addr:     addr,
		LastSeen: time.Now(),
	}
}

func (peer *Peer) Signal() {
	peer.LastSeen = time.Now()
}