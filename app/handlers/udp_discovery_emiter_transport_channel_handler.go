package handlers

import "peer-to-peer/app/peer_comunication"

type UDPDiscoveryEmitterTransportChannelHandler struct {}

func (u *UDPDiscoveryEmitterTransportChannelHandler) OnClose(channel peer_comunication.ITransportChannel) {
	//do nothing, this is a sender discovery channel
}

func (u *UDPDiscoveryEmitterTransportChannelHandler) OnOpen(channel peer_comunication.ITransportChannel) {
	//do nothing, this is a sender discovery channel
}

func (u *UDPDiscoveryEmitterTransportChannelHandler) OnMessage(channel peer_comunication.ITransportChannel, message peer_comunication.TransportMessage) error {
	return nil
}