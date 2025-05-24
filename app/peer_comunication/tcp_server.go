package peer_comunication

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

type TCPServer struct {
	Listener net.Listener
	channels map[TransportAddressKey]ITransportChannel
	port   int
}

func NewTCPServer(port int) *TCPServer {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Printf("Error starting TCP server on port %d: %v\n", port, err)
		panic(err)
	}

	return &TCPServer{
		Listener: listener,
		channels: make(map[TransportAddressKey]ITransportChannel),
		port:     port,
	}
}

func (s *TCPServer) Listen() error {
	log.Printf("TCP server listening on port %d\n", s.port)
	for {
		channel, err := s.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		addr := channel.GetAddress()
		log.Printf("Accepted connection from %s\n", addr.String())
	}
}

func (s *TCPServer) GetPort() int {
	return s.port
}


func (s *TCPServer) Accept() (ITransportChannel, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, err
	}

	channel := NewTCPTransportChannel(conn)
	address := TransportAddress{
		ip:   conn.RemoteAddr().(*net.TCPAddr).IP,
		port: conn.RemoteAddr().(*net.TCPAddr).Port,
	}
	s.channels[address.GetKey()] = channel

	RegisterTransportChannel(channel)

	go readMessageFromConn(channel, conn)

	return channel, nil
}

func (s *TCPServer) GetChannel(address TransportAddress) (ITransportChannel, bool) {
	key := address.GetKey()
	channel, exists := s.channels[key]
	return channel, exists
}

func (s *TCPServer) Close() error {
	if err := s.Listener.Close(); err != nil {
		return err
	}
	for _, channel := range s.channels {
		if err := channel.Close(); err != nil {
			return err
		}
	}
	s.channels = make(map[TransportAddressKey]ITransportChannel)
	return nil
}

func readMessageFromConn(channel ITransportChannel, conn net.Conn) error {
	for{
		// Read the size of the message
		sizeBytes := make([]byte, 4)
		_, err := conn.Read(sizeBytes)
		if err != nil {
			return err
		}
		size := binary.BigEndian.Uint32(sizeBytes)

		// Read the message
		messageBytes := make([]byte, size)
		_, err = conn.Read(messageBytes)
		if err != nil {
			return err
		}

		address := TransportAddress{
			ip:   conn.RemoteAddr().(*net.TCPAddr).IP,
			port: conn.RemoteAddr().(*net.TCPAddr).Port,
		}

		transportMessage := NewTransportMessage(size, messageBytes, address)
		err = channel.CollectMessage(transportMessage)
		if err != nil {
			log.Printf("Error collecting message from channel: %v\n", err)
			return err
		}
	}
}

