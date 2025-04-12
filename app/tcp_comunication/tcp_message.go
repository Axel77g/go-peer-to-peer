package tcpcomunication

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

type TCPMessage struct {
	size    uint32
	message []byte
	received bool
}

func NewTCPMessage(size uint32, message []byte, received bool) TCPMessage {
	return TCPMessage{
		size:    size,
		message: message,
		received: received,
	}
}

func ParseJSONMessage[T any](jsonMessage *TCPMessageJSON) (*T, error) {
	var payload T
	err := json.Unmarshal(jsonMessage.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, err
}

func (m *TCPMessage) Send(conn net.Conn) (bool, error) {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, m.size)
	_, err := conn.Write(lenBytes)
	if err != nil {
		return  false, err
	}
	//send the message
	_, err = conn.Write(m.message)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *TCPMessage) IsJSON() bool {
	if len(m.message) < 2 { return false }
	return m.message[0] == '{' && m.message[len(m.message)-1] == '}'
}

func (m *TCPMessage) GetJSON() (*TCPMessageJSON, error) {
	if !m.IsJSON() {
		return nil, nil
	}

	var messageJSON TCPMessageJSON
	err := json.Unmarshal(m.message, &messageJSON)
	if err != nil {
		return nil, err
	}

	return &messageJSON, nil
}

func ReceiveTCPMessage(conn net.Conn) (*TCPMessage, error) {
	lenBuf := make([]byte, 4) // 4 octets pour la taille du message
	_, err := conn.Read(lenBuf)
	if err != nil {
		return nil, err
	}
	msgLen := binary.BigEndian.Uint32(lenBuf)

	// Lire le message complet
	msgBuf := make([]byte, msgLen)
	_, err = conn.Read(msgBuf)
	if err != nil {
		return nil, err
	}

	message := NewTCPMessage(msgLen, msgBuf, true)
	return &message, nil
}

func CreateTCPMessageJSON(messageType string, payload any) (*TCPMessage, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	message := TCPMessageJSON{
		Type: messageType,
		Payload: jsonPayload,
	}
	
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	tcpMessage := TCPMessage{
		size: uint32(len(jsonMessage)),
		message: jsonMessage,
		received: false,
	}

	return &tcpMessage, nil
}

