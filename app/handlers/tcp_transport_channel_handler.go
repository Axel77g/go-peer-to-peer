package handlers

import (
	"log"
	"peer-to-peer/app/peer_comunication"
)

type TCPTransportChannelHandler struct {}

func (t *TCPTransportChannelHandler) OnClose(channel peer_comunication.ITransportChannel) {
	// Unregister the transport channel when it is closed
	peer_comunication.UnregisterTransportChannel(channel)
}

func (t *TCPTransportChannelHandler) OnOpen(channel peer_comunication.ITransportChannel) {
	// Register the transport channel when it is opened
	peer_comunication.RegisterTransportChannel(channel)
}

func (t *TCPTransportChannelHandler) OnMessage(channel peer_comunication.ITransportChannel, message peer_comunication.TransportMessage) error {
	log.Printf("Received message on TCP channel: %s\n", string(message.GetContent()))
	return nil
}