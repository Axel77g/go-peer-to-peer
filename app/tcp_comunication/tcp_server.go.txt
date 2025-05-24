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
	log.Println("[TCPServer] New connection from ", conn.RemoteAddr().String())
	remoteAddr := conn.RemoteAddr().String()
	
	t.mu.Lock()
	socket := TCPSocket{
		RemoteAddr: remoteAddr,
		Conn: conn,
		PeerID: "UNKOWN", // wait for the peer to be sent (wait for handshake)
	}
	t.sockets[remoteAddr] = socket	
	t.mu.Unlock()

	socket.ListenMessage(conn, t)
	t.removeSocket(socket)
}

func (t *TCPServer) removeSocket(socket TCPSocket) {
	log.Println("[TCPServer] Remove socket of peer: ", socket.PeerID)
	delete(t.sockets, socket.RemoteAddr)
}

func (t *TCPServer) PrintSocket(){
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, socket := range t.sockets {
		log.Printf("[TCP] Socket : %s, Addr : %s -> PeerID  %s \n ", socket.RemoteAddr, socket.Conn.RemoteAddr().String(), socket.PeerID)
	}
}