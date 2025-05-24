package peer_comunication

import (
	"errors"
	"log"
	"net"
)

type UDPTransportChannel struct {
	incoming chan ITransportMessage
	addr     net.Addr
	port     int
	stop     chan struct{}
}

func NewUDPTransportChannel(port int) *UDPTransportChannel {
	channel := &UDPTransportChannel{
		port:     port,
		incoming: make(chan ITransportMessage, 100),
		stop:     make(chan struct{}),
	}
	go channel.readLoop()
	return channel
}

func (u *UDPTransportChannel) readMessage(conn net.UDPConn) {
	messageSize := make([]byte, 4)
	size, remoteAddr := conn.ReadFromUDP(messageSize)

}

func (u *UDPTransportChannel) readLoop() {
	// open udp server on port
	// listen for incoming messages

	addr := net.UDPAddr{
		Port: u.port,
		IP:   net.IPv4zero,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Printf("Erreur lors de l'écoute UDP : %v\n", err)
		return
	}

	defer conn.Close()
	log.Println("UDP server listening on port", u.port)

	for {
		select {
		case <-u.stop:
			return
		default:
			message, err := u.readMessage()
			if err != nil {
				continue
			}
			u.incoming <- message
		}
	}
	/* messageSize := make([]byte, 4)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Erreur de lecture UDP : %v\n", err)
			continue
		}

		// process incomming message
		// the transport channels are use for remote peer communication and also for the local peer to write on network

	} */

}

func (u *UDPTransportChannel) GetPort() int {
	return u.port
}

func (u *UDPTransportChannel) Send(content []byte) error {
	// Implement UDP send logic here
	return nil
}
func (u *UDPTransportChannel) Read() (ITransportMessage, error) {
	// Implement UDP read logic here
	return nil, nil
}

func (u *UDPTransportChannel) Close() error {
	return errors.New("Cannot close UDP transport channel")
}
