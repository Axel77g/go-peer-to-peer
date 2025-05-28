package handlers

import (
	"log"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
)

type TCPControllerTransportChannelHandler struct{}

func (t *TCPControllerTransportChannelHandler) OnClose(channel peer_comunication.ITransportChannel) {
	// Unregister the transport channel when it is closed
	peer_comunication.UnregisterTransportChannel(channel)
}

func (t *TCPControllerTransportChannelHandler) OnOpen(channel peer_comunication.ITransportChannel) {
	// Register the transport channel when it is opened
	peer_comunication.RegisterTransportChannel(channel)

	address := channel.GetAddress()
	isClientChannel := address.GetPort() == shared.TCPPort // si le port de destination est le port du serveur, c'est un channel client
	if isClientChannel {                                   // on pull au connection si on est client
		t.pullRemoteEvents(channel)
	}
	//on connect in tcp ask pull to remote peer is directory shallow

}

func (t *TCPControllerTransportChannelHandler) OnMessage(channel peer_comunication.ITransportChannel, message peer_comunication.TransportMessage) error {
	log.Printf("Received message on TCP channel: %s\n", string(message.GetContent()))

	if string(message.GetContent()) == "PULL_EVENTS_REQUEST" {
		/* collection := file_event.NewJSONLFileEventCollection("events.jsonl") */
		/* size := collection.GetBytesSize()
		iterator := collection.GetIterator()


		*/
	}

	return nil
}

func (t *TCPControllerTransportChannelHandler) pullRemoteEvents(channel peer_comunication.ITransportChannel) {
	channel.Send([]byte("PULL_EVENTS_REQUEST"))
}
