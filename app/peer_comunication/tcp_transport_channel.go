package peer_comunication

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	file_event "peer-to-peer/app/files/event"
	"peer-to-peer/app/shared"
	"time"
)

type TCPTransportChannel struct {
	conn            net.Conn
	lastMessageTime time.Time
	handler         ITransportChannelHandler // Handler for transport channel events
}

func NewTCPTransportChannel(conn net.Conn, handler ITransportChannelHandler) *TCPTransportChannel {
	t := &TCPTransportChannel{
		conn:    conn,
		handler: handler,
	}
	handler.OnOpen(t) // Notify the handler that the channel is opened
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

func (t *TCPTransportChannel) SendIterator(size uint32, message []byte, iterator shared.Iterator) error {
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, size + uint32(len(message)))
	_, err := t.conn.Write(sizeBytes)
	if err != nil {
		return err
	}

	// Send the initial message
	_, err = t.conn.Write(message)

	for iterator.Next() {
		content, err := iterator.Current()
		if err != nil {
			log.Printf("Error getting current content from iterator: %v\n", err)
			panic("Failed to get current content from iterator: " + err.Error())
		}
		
		// Convert any to bytes
		var contentBytes []byte
		switch v := content.(type) {
		case []byte:
			contentBytes = v
		case string:
			contentBytes = []byte(v)
		case file_event.FileEvent:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				log.Printf("Error marshaling FileEvent to JSON: %v\n", err)
				contentBytes = []byte(fmt.Sprintf("%v", v))
			} else {
				contentBytes = jsonBytes
			}
		default:
			contentBytes = []byte(fmt.Sprintf("%v", v))
		}
		
		_, err = t.conn.Write(contentBytes)
		if err != nil {
			panic("Failed to send content: " + err.Error())
		}
	}
	return nil
}

func (t *TCPTransportChannel) CollectMessage(message TransportMessage) error {
	t.lastMessageTime = time.Now()
	t.handler.OnMessage(t, message) // Notify the handler that a message has been collected
	return nil
}

func (t *TCPTransportChannel) Close() error {
	t.handler.OnClose(t) // Notify the handler that the channel is closed
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
	for {
		// Read the size of the message
		sizeBytes := make([]byte, 4)
		_, err := t.conn.Read(sizeBytes)
		if err != nil {
			if err.Error() == "EOF" || err.Error() == "use of closed network connection" {
				t.Close()
			}
			log.Printf("Error reading TCP buffer from connection: %v\n", err)
			return err
		}
		size := binary.BigEndian.Uint32(sizeBytes)

		// Read the message
		messageBytes := make([]byte, size)
		_, err = t.conn.Read(messageBytes)
		if err != nil {
			if err.Error() == "EOF" || err.Error() == "use of closed network connection" {
				t.Close()
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