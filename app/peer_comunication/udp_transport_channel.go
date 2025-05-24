package peer_comunication

import (
	"encoding/binary"
	"log"
	"net"
)


type UDPTransportChannel struct {
	listener *UDPServerListener
	address TransportAddress
	incoming chan TransportMessage
}

func NewUDPTransportChannel(address TransportAddress) *UDPTransportChannel {
	log.Printf("New UDP transport channel created for address: %s\n", address.String())
	listener := GetUDPServerListener() //get the udp listener singleton
	return &UDPTransportChannel{
		listener: listener,
		address:  address,
	}
}

func (u *UDPTransportChannel) GetPort() int {
	return u.address.port
}

func (u *UDPTransportChannel) GetAddress() TransportAddress {
	return u.address
}

func (u *UDPTransportChannel) Send(content []byte) error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   u.address.ip,
		Port: u.GetPort(),
	})

	if err != nil {
		return err
	}

	defer conn.Close()

	size := uint32(len(content))
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, size)
	_, err = conn.Write(sizeBytes)
	if err != nil {
		return err
	}

	_, err = conn.Write(content)
	if err != nil {
		return err
	}
	return nil
}


func (u *UDPTransportChannel) Close() error {
	// UDP does not have a close method like TCP, but we can return nil
	return nil
}

func (u *UDPTransportChannel) Read() (TransportMessage, error) {
	select {
	case message := <-u.incoming:
		return message, nil
	default:
		return TransportMessage{}, nil // No message available
	}
}

func (u *UDPTransportChannel) CollectMessage(message TransportMessage) error {
	u.incoming <- message
	return nil
}

