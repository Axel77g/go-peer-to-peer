package tcpcomunication

import (
	"log"
	"net"
)

type TCPServerSocket struct {
	Identifier string
	Conn net.Conn
	PeerID string
}

func (socket *TCPServerSocket) ListenForMessage(ts *TCPServer, conn net.Conn) {
	// Handle the message
	for {
		message, err := ReceiveTCPMessage(conn)
		if err != nil{
			log.Fatalln("Error while receiving message: ", err)
			break
		}
		
		if message == nil {
			log.Println("Message is nil")
			break
		}

		if message.IsJSON() {
			jsonMessage, err := message.GetJSON()
			if err != nil {
				log.Println("Error while getting JSON message: ", err)
				break
			}

			if(jsonMessage == nil) {
				log.Println("JSON message is nil")
				break
			}

			ts.mu.Lock()
			if !socket.HandshakeEnd() && jsonMessage.Type != MESSAGE_TYPE_HELLO {
				log.Println("Handshake not finished, ignoring message")
				continue
			}

			log.Println("Received message: ", jsonMessage)

			switch jsonMessage.Type {
				case MESSAGE_TYPE_HELLO:
					hello, _ := ParseJSONMessage[HelloMessage](jsonMessage)
					socket.setPeerID(hello.PeerID)
					ts.sockets[socket.Identifier] = *socket
				case MESSAGE_TYPE_FILE_DIR:
					dir, _ := ParseJSONMessage[FileDirMessage](jsonMessage)
					log.Println("Received directory: ", dir.Files)
					//compare the directory with the one in the peer
			}

			ts.mu.Unlock()
		}
    }
}

func (t *TCPServerSocket) setPeerID(peerID string) {
	t.PeerID = peerID
}

func (t *TCPServerSocket) HandshakeEnd() bool {
	return t.PeerID != ""
}
