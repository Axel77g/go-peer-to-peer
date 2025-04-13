package tcpcomunication

import (
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
	PeerID string // can be nil in the beginning (handshake not started and not finished), it must be a reference of peerManager peer
}

func (socket *TCPSocket) ListenMessage(conn net.Conn, mu *sync.Mutex) {
	defer socket.HandleDeconnection()
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
					peerInstance,exist := peerManager.GetPeer(hello.PeerID)
					if !exist {
						log.Println("Receive HELLO from a peer that is not in the peer manager :", hello.PeerID)
						peerInstance = peer.NewPeer(hello.PeerID, conn.RemoteAddr().(*net.TCPAddr).IP)
					}else{
						log.Println("Receive HELLO from a peer that is already in the peer manager :", hello.PeerID)
					}
					peerInstance.TCPSocket = socket
					socket.setPeerID(peerInstance.ID)
					peerManager.UpsertPeer(peerInstance)
				case MESSAGE_TYPE_FILE_DIR:
					dir, _ := ParseJSONMessage[FileDirMessage](jsonMessage)
					log.Println("Received directory: ", dir.Files)
					//compare the directory with the one in the peer
			}

			mu.Unlock()
		}
    }
}

func (socket *TCPSocket) HandleDeconnection() {
	if socket.Conn != nil {
		socket.Conn.Close()
	}
	socket.Conn = nil
	peerManager := peer.GetPeerManager()
	peer, exist := peerManager.GetPeer(socket.PeerID)
	if exist {
		peer.TCPSocket = nil
		peerManager.UpsertPeer(peer)
		log.Println("Peer disconnected: ", socket.PeerID)
	}
}

func (t *TCPSocket) Close() {
	if t.Conn != nil {
		t.Conn.Close()
	}
	t.Conn = nil
}

func (t *TCPSocket) Send(message TCPMessage) (bool, error) {
	if t.Conn == nil {
		return false, nil
	}
	return message.Send(t.Conn)
}

func (t *TCPSocket) setPeerID(peer string) {
	t.PeerID = peer
}

func (t *TCPSocket) HandshakeEnd() bool {
	return t.PeerID != "UNKOWN"
}


func CreateTCPConnection(peer *peer.Peer) (TCPSocket,error) {
	log.Println("Creating TCP (client) connection to peer: ", peer.ID)
	conn, err := net.Dial("tcp", peer.Addr.String() + ":" + strconv.Itoa(shared.TCPPort))
	if err != nil {
		return TCPSocket{}, err
	}
	socket := TCPSocket{
		RemoteAddr: conn.RemoteAddr().String(),
		Conn: conn,
		PeerID: peer.ID,
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