package tcpcomunication

import (
	"errors"
	"io"
	"log"
	"net" // Renamed to avoid confusion, but no structural change here.
	"peer-to-peer/app/peer"
	"peer-to-peer/app/shared"
	"strconv"
	"sync"
)

// interface TCPsocket
type TCPSocket struct {
	RemoteAddr string
	Conn net.Conn
	Peer *peer.Peer // can be nil in the beginning (handshake not started and not finished), it must be a reference of peerManager peer
}

func (socket *TCPSocket) ListenMessage(conn net.Conn, mu *sync.Mutex) {
	// Handle the message
	for {
		message, err := ReceiveTCPMessage(conn)
		if err != nil{
			if err == io.EOF{
				log.Println("Read end connection closed", socket.RemoteAddr)
				break
			}

			log.Println("Error while receiving message: ", err)
			break
		}
	

		if message.IsJSON() {
			jsonMessage, err := message.GetJSON()
			if err != nil {
				log.Println("Error while getting JSON message: ", err)
				break
			}

			mu.Lock()
			if !socket.HandshakeEnd() && jsonMessage.Type != MESSAGE_TYPE_HELLO {
				log.Println("Handshake not finished, ignoring message")
				continue
			}

			log.Println("Received message: ", jsonMessage)

			switch jsonMessage.Type {
				case MESSAGE_TYPE_HELLO:
					hello, _ := ParseJSONMessage[HelloMessage](jsonMessage)
					peerManager := peer.GetPeerManager()
					peerRef := peerManager.GetPeer(hello.PeerID)
					if peerRef == nil {
						log.Println("Peer not found, creating new peer")
						newPeer := peer.NewPeer(hello.PeerID, conn.RemoteAddr().(*net.TCPAddr).IP)
						peerManager.SignalPeer(newPeer) //signal only if the peer is new add it to the peer manager
						peerRef = peerManager.GetPeer(hello.PeerID)
					}
					socket.setPeer(peerRef)
				case MESSAGE_TYPE_FILE_DIR:
					dir, _ := ParseJSONMessage[FileDirMessage](jsonMessage)
					log.Println("Received directory: ", dir.Files)
					//compare the directory with the one in the peer
			}

			mu.Unlock()
		}
    }
}

func (t *TCPSocket) Send(message TCPMessage) (bool, error) {
	if t.Conn == nil {
		return false, nil
	}
	return message.Send(t.Conn)
}

func (t *TCPSocket) GetConn() (net.Conn, error) {
	if t.Conn == nil {
		return nil, errors.New("socket not connected")
	}
	return t.Conn, nil
}

func (t *TCPSocket) setPeer(peer *peer.Peer) {
	t.Peer = peer
}

func (t *TCPSocket) HandshakeEnd() bool {
	return t.Peer != nil
}


func CreateTCPConnection(peer peer.Peer) (TCPSocket,error) {
	conn, err := net.Dial("tcp", peer.Addr.String() + ":" + strconv.Itoa(shared.TCPPort))
	if err != nil {
		return TCPSocket{}, err
	}
	socket := TCPSocket{
		RemoteAddr: conn.RemoteAddr().String(),
		Conn: conn,
		Peer: &peer,
	}
	helloMessage := HelloMessage{
		PeerID: strconv.Itoa(shared.SOCKET_ID),
	}
	message, err := CreateTCPMessageJSON(MESSAGE_TYPE_HELLO, helloMessage)
	if err != nil {
		return socket, err
	}

	if _, err := socket.Send(message); err != nil {
		return socket, err
	}
	go socket.ListenMessage(conn, &sync.Mutex{})
	return socket, nil
}