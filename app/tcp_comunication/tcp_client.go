package tcpcomunication

import (
	"net"
	"peer-to-peer/app/peer"
	"peer-to-peer/app/shared"
	"strconv"
)

type TCPClient struct {
	Conn net.Conn
	Peer peer.Peer
}

func NewTCPClient(peer peer.Peer) TCPClient {
	return TCPClient{
		Conn: nil,
		Peer: peer,
	}
}

func (t *TCPClient) Connect() error {
	conn, err := net.Dial("tcp", t.Peer.Addr.String() + ":" + strconv.Itoa(shared.TCPPort))
	if err != nil {
		return err
	}
	t.Conn = conn
	helloMessage := HelloMessage{
		PeerID: t.Peer.ID,
	}
	message, err := CreateTCPMessageJSON(MESSAGE_TYPE_HELLO, helloMessage)
	if err != nil {
		return err
	}

	if _, err := message.Send(t.Conn); err != nil {
		return err
	}
	defer t.Conn.Close()
	return nil
}