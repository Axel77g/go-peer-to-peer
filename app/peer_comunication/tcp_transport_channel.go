package peer_comunication

import (
	"encoding/binary"
	"net"
)

type TCPTransportChannel struct {
	conn     net.Conn
	incoming chan ITransportMessage
	stop     chan struct{}
}

func NewTcpTransportChannel(conn net.Conn) *TCPTransportChannel {
	t := &TCPTransportChannel{
		conn:     conn,
		incoming: make(chan ITransportMessage, 100),
		stop:     make(chan struct{}),
	}
	go t.readLoop()
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

func (t *TCPTransportChannel) readMessage() (ITransportMessage, error) {
	// Read the size of the message
	sizeBytes := make([]byte, 4)
	_, err := t.conn.Read(sizeBytes)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(sizeBytes)

	// Read the message
	messageBytes := make([]byte, size)
	_, err = t.conn.Read(messageBytes)
	if err != nil {
		return nil, err
	}

	return NewTCPMessage(size, messageBytes, false), nil
}

func (t *TCPTransportChannel) readLoop() {
	for {
		select {
		case <-t.stop:
			return
		default:
			message, err := t.readMessage()
			if err != nil {
				continue
			}
			t.incoming <- message
		}
	}
}

func (t *TCPTransportChannel) Read() (ITransportMessage, error) {
	select {
	case message := <-t.incoming:
		return message, nil
	case <-t.stop:
		return nil, nil
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