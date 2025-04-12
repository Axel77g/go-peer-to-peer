package peer

import "time"

type Peer struct {
	ID       string
	Addr     string
	TCPPort  uint16
	LastSeen time.Time
}

func NewPeer(id, addr string, tcpPort uint16) Peer {
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