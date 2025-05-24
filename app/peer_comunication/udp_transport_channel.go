package peer_comunication

import (
	"errors"
	"net"
)


type UDPTransportChannel struct {
	listener *UDPServerListener
	address TransportAddress
	incoming chan TransportMessage
}

func NewUDPTransportChannel(address TransportAddress) *UDPTransportChannel {
	listener := GetUDPServerListener() //get the udp listener singleton
	return &UDPTransportChannel{
		listener: listener,
		address:  address,
		incoming: make(chan TransportMessage, 100), // Buffered channel for incoming messages
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
		panic("Failed to connect to UDP server: " + err.Error())
	}

	defer conn.Close()

	_, err = conn.Write(content)
	if err != nil {
		panic("Failed to send content: " + err.Error())
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
	select {
		case u.incoming <- message:
		default:
			return errors.New("channel full") // If the channel is full, return an error
	}
	return nil
}
