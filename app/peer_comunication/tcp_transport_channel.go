package peer_comunication

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
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
	println("Sending message of size on TCP socket:", len(content))
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

func (t *TCPTransportChannel) SendIterator(message []byte, iterator shared.Iterator) error {
	// Calculer la taille des données (sans inclure les 4 bytes d'en-tête)
	messageSize := uint32(len(message))
	iteratorSize := uint32(0)
	// Calculer la taille totale en parcourant l'iterator une première fois
	var contents [][]byte
	for iterator.Next() {
		content, err := iterator.Current()
		if err != nil {
			log.Printf("Error getting current content from iterator: %v\n", err)
			return fmt.Errorf("failed to get current content from iterator: %w", err)
		}

		// Convert any to bytes
		var contentBytes []byte
		switch v := content.(type) {
		case []byte:
			contentBytes = v
		case string:
			contentBytes = []byte(v)
		case shared.FileEvent:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				log.Printf("Error marshaling FileEvent to JSON: %v\n", err)
				contentBytes = []byte(fmt.Sprintf("%v\n", v))
			} else {
				contentBytes = append(jsonBytes, "\n"...)
			}
		default:
			contentBytes = []byte(fmt.Sprintf("%v", v))
		}

		contents = append(contents, contentBytes)
		iteratorSize += uint32(len(contentBytes))
	}

	dataSize := messageSize + iteratorSize
	// Taille totale = taille des données + 4 bytes pour l'en-tête
	totalSize := 4 + dataSize
	buffer := make([]byte, 0, totalSize)

	// Set la taille des données dans les 4 premiers bytes
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, dataSize)
	buffer = append(buffer, sizeBytes...)

	// Ajouter le message initial
	buffer = append(buffer, message...)

	// Ajouter tout le contenu de l'iterator
	for _, contentBytes := range contents {
		buffer = append(buffer, contentBytes...)
	}

	// Envoyer tout en un seul Write
	_, err := t.conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("failed to send content: %w", err)
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
		_, err := io.ReadFull(t.conn, sizeBytes)
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
		_, err = io.ReadFull(t.conn, messageBytes)
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
