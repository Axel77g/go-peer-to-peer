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

func NewPeer(id string, addr net.IP, tcpPort uint16) Peer {
	return Peer{
		ID:       id,
		Addr:     addr,
		TCPPort:  tcpPort,
		LastSeen: time.Now(),
	}
}

func (peer *Peer) setTCPPort(port uint16) {
	peer.TCPPort = port
}