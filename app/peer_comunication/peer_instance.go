package peer_comunication

import (
	"net"
	"strconv"
)

type IPeer interface {
	GetIP() net.IP
	GetTransportsChannels() *TransportChannelCollection
	AddTransportChannel(channel ITransportChannel)
	RemoveTransportChannel(channel ITransportChannel)
	String() string
}

type Peer struct {
	address             net.IP
	transportsChannels  *TransportChannelCollection
}

func NewPeer(address net.IP) *Peer {
	return &Peer{
		address: address,
		transportsChannels: &TransportChannelCollection{
			channels: make([]ITransportChannel, 0),
		},
	}
}

func (p *Peer) GetIP() net.IP {
	return p.address
}

func (p *Peer) GetTransportsChannels() *TransportChannelCollection {
	return p.transportsChannels
}

func (p *Peer) AddTransportChannel(channel ITransportChannel) {
	p.transportsChannels.Add(channel)
}

func (p *Peer) RemoveTransportChannel(channel ITransportChannel) {
	p.transportsChannels.Remove(channel)
}

func (p *Peer) String() string {
	return p.address.String() + " - " + strconv.Itoa(p.transportsChannels.Size()) + " channels"
}