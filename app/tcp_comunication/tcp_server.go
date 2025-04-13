package tcpcomunication

import (
	"log"
	"net"
	"peer-to-peer/app/shared"
	"strconv"
	"sync"
)

type TCPServer struct  {
	Listener net.Listener
	sockets map[string]TCPSocket
	mu sync.Mutex
}

func NewTCPServer() *TCPServer {
	return &TCPServer{
		Listener: nil,
		sockets: make(map[string]TCPSocket),
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
	
	log.Println("[TCPServer] New connection from ", conn.RemoteAddr().String())
	remoteAddr := conn.RemoteAddr().String()
	
	t.mu.Lock()
	tcpPeerConn := TCPSocket{
		RemoteAddr: remoteAddr,
		Conn: conn,
		Peer: nil, // wait for the peer to be sent (wait for handshake)
	}
	t.sockets[remoteAddr] = tcpPeerConn	
	t.mu.Unlock()

	tcpPeerConn.ListenMessage(conn, &t.mu)
}

func (t *TCPServer) handleDeconnection(conn net.Conn) {
	log.Println("[TCPServer] Socket deconnection")
	remoteAddr := conn.RemoteAddr().String()
	delete(t.sockets, remoteAddr)
	conn.Close()
}