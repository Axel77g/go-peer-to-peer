package peer_comunication

import (
	"net"
	"sync"
	"time"
)


type UDPTransportChannel struct {
	listener *UDPServerListener
	address TransportAddress
	incoming chan TransportMessage
	stop    chan struct{}
	lastMessageTime time.Time
	monitorOnce      sync.Once
}

func NewUDPTransportChannel(address TransportAddress) *UDPTransportChannel {
	listener := GetUDPServerListener() //get the udp listener singleton
	return &UDPTransportChannel{
		listener: listener,
		address:  address,
		incoming: make(chan TransportMessage, 100), // Buffered channel for incoming messages
		stop:     make(chan struct{}),
		lastMessageTime: time.Now(),
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
	// Block until a message is available or the channel is closed
	select {
		case message := <-u.incoming:
			return message, nil
		case <-u.stop:
			return TransportMessage{}, nil // Channel closed, no message available
	}
}

func (u *UDPTransportChannel) CollectMessage(message TransportMessage) error {
	select {
		case u.incoming <- message:
		default:
			<- u.incoming // If the channel is full, drop the oldest message
			u.incoming <- message // and add the new message
	}
	u.lastMessageTime = time.Now()
	u.monitorOnce.Do(func() {
		go func() {
			time.Sleep(20 * time.Second)
			if !u.IsAlive() {
				UnregisterTransportChannel(u) // Unregister the channel if it is not alive
			}
		}()
	})
	return nil
}


func (u *UDPTransportChannel) GetProtocol() string {
	return "udp"
}

func (u *UDPTransportChannel) IsAlive() bool {
	return time.Since(u.lastMessageTime) < 20*time.Second // Consider alive if last message was received within 30 seconds
}
