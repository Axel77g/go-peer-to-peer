package tcpcomunication

import (
	"net"
	"peer-to-peer/app/shared"
	"strconv"
	"sync"
)

type TCPServer struct {
	Listener net.Listener
	sockets map[string]TCPServerSocket
	mu sync.Mutex
}

func NewTCPServer() *TCPServer {
	return &TCPServer{
		Listener: nil,
		sockets: make(map[string]TCPServerSocket),
		mu: sync.Mutex{},
	}
}

// go routine to listen for incoming connections
func (t *TCPServer) Listen() {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(shared.TCPPort))
	if err != nil {
		panic(err)
	}
	t.Listener = listener

	defer t.Listener.Close()
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			continue
		}
		go t.handleConnection(conn)
	}
}

// go routine to handle the connection
func (t *TCPServer) handleConnection(conn net.Conn) {
	defer t.handleDeconnection(conn)
	
	identifier := conn.RemoteAddr().String()
	
	t.mu.Lock()
	tcpPeerConn := TCPServerSocket{
		Conn: conn,
		PeerID: "", // wait for the peer ID to be sent (handshake)
	}
	t.sockets[identifier] = tcpPeerConn	
	t.mu.Unlock()

	tcpPeerConn.ListenForMessage(t, conn)
}

func (t *TCPServer) handleDeconnection(conn net.Conn) {
	identifier := conn.RemoteAddr().String()
	delete(t.sockets, identifier)
	conn.Close()
}