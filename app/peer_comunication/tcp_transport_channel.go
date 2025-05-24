package peer_comunication

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

type TCPTransportChannel struct {
	conn     net.Conn
	incoming chan TransportMessage
	stop     chan struct{}
	lastMessageTime time.Time
}

func NewTCPTransportChannel(conn net.Conn) *TCPTransportChannel {
	t :=  &TCPTransportChannel{
		conn:     conn,
		incoming: make(chan TransportMessage, 100),
		stop:     make(chan struct{}),
	}
	go t.readMessageFromConn()
	return t
}

func (t *TCPTransportChannel) Send(content []byte) error {
	// Send the size of the message
	size := uint32(len(content))
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, size)
	_, err := t.conn.Write(sizeBytes)
	if err != nil {
		return err
	}

	// Send the message
	_, err = t.conn.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func (t *TCPTransportChannel) CollectMessage(message TransportMessage) error {
	select {
		case t.incoming <- message:
		default:
			return errors.New("channel full") // If the channel is full, return an error
	}
	t.lastMessageTime = time.Now()
	return nil
}

func (t *TCPTransportChannel) Read() (TransportMessage, error) {
	select {
		case message := <-t.incoming:
			return message, nil
		case <-t.stop:
			return TransportMessage{}, nil
	}
}

func (t *TCPTransportChannel) Close() error {
	close(t.stop)
	return t.conn.Close()
}

func (t *TCPTransportChannel) GetPort() int {
	if tcpAddr, ok := t.conn.RemoteAddr().(*net.TCPAddr); ok {
		return tcpAddr.Port
	}
	return 0
}

func (t *TCPTransportChannel) GetAddress() TransportAddress {
	if tcpAddr, ok := t.conn.RemoteAddr().(*net.TCPAddr); ok {
		return TransportAddress{
			ip:   tcpAddr.IP,
			port: tcpAddr.Port,
		}
	}
	return TransportAddress{}
}

func (t *TCPTransportChannel) GetProtocol() string {
	return "tcp"
}

func (t *TCPTransportChannel) IsAlive() bool {
	return t.conn != nil && t.lastMessageTime.Add(10*time.Minute).After(time.Now())
}


func (t *TCPTransportChannel) readMessageFromConn() error {
	for{
		// Read the size of the message
		sizeBytes := make([]byte, 4)
		_, err := t.conn.Read(sizeBytes)
		if err != nil {
			if(err.Error() == "EOF" || err.Error() == "use of closed network connection") {
				handleDisconnection(t)
			}
			log.Printf("Error reading TCP buffer from connection: %v\n", err)
			return err
		}
		size := binary.BigEndian.Uint32(sizeBytes)

		// Read the message
		messageBytes := make([]byte, size)
		_, err = t.conn.Read(messageBytes)
		if err != nil {
			if(err.Error() == "EOF" || err.Error() == "use of closed network connection") {
				handleDisconnection(t)
			}
			log.Printf("Error reading TCP buffer from connection: %v\n", err)
			return err
		}

	
		transportMessage := NewTransportMessage(size, messageBytes, t.GetAddress())
		err = t.CollectMessage(transportMessage)
		if err != nil {
			log.Printf("Error collecting message from channel: %v\n", err)
			return err
		}
	}
}

func handleDisconnection(channel ITransportChannel) {
	address := channel.GetAddress()
	log.Printf("Handling disconnection for address: %s\n", address.String())
	UnregisterTransportChannel(channel)
	if err := channel.Close(); err != nil {
		log.Printf("Error closing channel: %v\n", err)
	}
}