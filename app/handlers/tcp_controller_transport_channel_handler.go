package handlers

import (
	"log"
	file_event "peer-to-peer/app/files/event"
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
	//log.Printf("Received message on TCP channel: %s\n", string(message.GetContent()))
	log.Printf("Received message size: %d bytes\n", message.GetSize())

	content := message.GetContent()
	stringContent := string(content)

	eventManager := file_event.GetEventManager()

	if stringContent == "PULL_EVENTS_REQUEST" {
		eventManager.Lock()
		collection := eventManager.GetCollection()
		size := collection.GetBytesSize()
		log.Printf("Sending events to remote peer, size: %d bytes\n", size)
		iterator := collection.GetAll("pulling events")
		defer iterator.Close()
		if iterator == nil {
			log.Println("No events found to send.")
			return nil
		}
		adapter := file_event.NewFileEventIteratorAdapter(iterator)
		messageContent := []byte("PULL_EVENTS_RESPONSE")
		channel.SendIterator(messageContent, adapter)
		eventManager.Unlock()
		log.Printf("Sent events to remote peer, size: %d bytes\n", size)
		return nil
	}

	//check the first parts of the message check if it PULL_EVENTS_RESPONSE, if case take the rest of the message and save it in the file events.jsonl
	pullEventResponseLen := len("PULL_EVENTS_RESPONSE")
	if len(content) > pullEventResponseLen && stringContent[:pullEventResponseLen] == "PULL_EVENTS_RESPONSE" {
		eventsData := content[pullEventResponseLen:]
		log.Printf("Received message PULL_EVENTS_RESPONSE, size: %d bytes\n", len(eventsData))

		address := channel.GetAddress()

		remote_collection := file_event.NewJSONLFileEventCollection("pull_events_from_remote_"+address.String()+".jsonl", true)
		err := remote_collection.FromBytes(eventsData)
		if err != nil {
			log.Printf("Error saving events from remote: %v\n", err)
			return err
		}
		log.Printf("Events from remote saved successfully in events_from_remote.jsonl")

		hasChanges, err := eventManager.MergeAndSave(remote_collection)
		if err != nil {
			log.Printf("Error merging and saving events: %v\n", err)
			eventManager.Unlock()
			return err
		}

		if hasChanges {
			log.Println("Changes detected, broadcasting events.")
			eventManager.BroadcastEvents()
		}

		return nil
	}

	/**
		Someone push events to us, we merge it with our events and save it if change we broadcast it again
	**/
	pushEventResponseLen := len("PUSH_EVENTS")
	if len(content) > pushEventResponseLen && stringContent[:pushEventResponseLen] == "PUSH_EVENTS" {
		eventsData := content[pushEventResponseLen:]
		log.Printf("Received message PUSH_EVENTS, size: %d bytes\n", len(eventsData))

		address := channel.GetAddress()
		remote_collection := file_event.NewJSONLFileEventCollection("push_events_to_remote_"+address.String()+".jsonl", true)
		err := remote_collection.FromBytes(eventsData)
		if err != nil {
			log.Printf("Error saving events to remote: %v\n", err)
			return err
		}
		log.Printf("Events to remote saved successfully in events_to_remote.jsonl")

		hasChanges, err := eventManager.MergeAndSave(remote_collection)
		if err != nil {
			log.Printf("Error merging and saving events: %v\n", err)
			eventManager.Unlock()
			return err
		}

		if hasChanges {
			log.Println("Changes detected, broadcasting events.")
			eventManager.BroadcastEvents()
		}

		return nil
	}

	return nil
}

func (t *TCPControllerTransportChannelHandler) pullRemoteEvents(channel peer_comunication.ITransportChannel) {
	channel.Send([]byte("PULL_EVENTS_REQUEST"))
}
