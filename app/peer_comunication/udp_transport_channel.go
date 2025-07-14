package peer_comunication

import (
	"log"
	"net"
	"peer-to-peer/app/shared"
	"sync"
	"time"
)

type UDPTransportChannel struct {
	listener        *UDPServerListener
	address         TransportAddress
	lastMessageTime time.Time
	monitorOnce     sync.Once
	eventHandler    ITransportChannelHandler // Handler for transport channel events
}

func NewUDPTransportChannel(address TransportAddress, handler ITransportChannelHandler) *UDPTransportChannel {
	listener := GetUDPServerListener() //get the udp listener singleton
	channel := &UDPTransportChannel{
		listener:        listener,
		address:         address,
		lastMessageTime: time.Now(),
		eventHandler:    handler,
	}
	handler.OnOpen(channel) // Notify the handler that the channel is opened
	return channel
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
	u.eventHandler.OnClose(u)
	return nil
}

func (u *UDPTransportChannel) CollectMessage(message TransportMessage) error {
	u.lastMessageTime = time.Now()
	u.monitorOnce.Do(func() {
		go func() {
			time.Sleep(10 * time.Second)
			log.Printf("Checking if UDP transport channel %s is alive\n", u.address.String())
			if !u.IsAlive() {
				u.Close()
				UnregisterTransportChannel(u) // Unregister the channel if it is not alive
			}
		}()
	})
	u.eventHandler.OnMessage(u, message)
	return nil
}

func (u *UDPTransportChannel) GetProtocol() string {
	return "udp"
}

func (u *UDPTransportChannel) IsAlive() bool {
	log.Printf("Time since last message: %v", time.Since(u.lastMessageTime))
	return time.Since(u.lastMessageTime) < 8*time.Second // Consider alive if last message was received within 30 seconds
}

func (u *UDPTransportChannel) SendIterator(message []byte, iterator shared.Iterator) error {
	// UDP doesn't support iterators in the same way as TCP
	// This is a simplified implementation - you might want to implement chunking
	log.Printf("Warning: SendIterator on UDP is not fully implemented")
	return nil
}
