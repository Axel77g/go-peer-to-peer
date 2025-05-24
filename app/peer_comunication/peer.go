package peer_comunication

import (
	"net"
	"strconv"
	"sync"
)

type IPeer interface {
	getAddress() net.IP
	getTransportsChannels() []ITransportChannel
	addTransportChannel(channel ITransportChannel)
	removeTransportChannel(channel ITransportChannel)
	String() string
}

type Peer struct {
	address             net.IP
	transportsChannel []ITransportChannel
	mutex             sync.RWMutex
}

func NewPeer(address net.IP) *Peer {
	return &Peer{
		address:             address,
		transportsChannel: []ITransportChannel{},
	}
}

func (p *Peer) getAddress() net.IP {
	return p.address
}

func (p *Peer) getTransportsChannels() []ITransportChannel {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.transportsChannel
}

func (p *Peer) addTransportChannel(channel ITransportChannel) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.transportsChannel = append(p.transportsChannel, channel)
}

func (p *Peer) removeTransportChannel(channel ITransportChannel) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	transportAddress := channel.GetAddress()
	for i, ch := range p.transportsChannel {
		chAddr := ch.GetAddress()
		if chAddr.GetKey() == transportAddress.GetKey() {
			p.transportsChannel = append(p.transportsChannel[:i], p.transportsChannel[i+1:]...)
			break
		}
	}
}

func (p *Peer) String() string {
	return p.address.String() + " - " + strconv.Itoa(len(p.transportsChannel)) + " channels"
}