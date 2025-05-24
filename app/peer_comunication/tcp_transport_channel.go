package peer_comunication

import (
	"encoding/binary"
	"net"
)

type TCPTransportChannel struct {
	conn     net.Conn
	incoming chan TransportMessage
	stop     chan struct{}
}

func NewTCPTransportChannel(conn net.Conn) *TCPTransportChannel {
	return &TCPTransportChannel{
		conn:     conn,
		incoming: make(chan TransportMessage, 100),
		stop:     make(chan struct{}),
	}
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
	t.incoming <- message
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