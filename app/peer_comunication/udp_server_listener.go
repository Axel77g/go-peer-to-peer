package peer_comunication

import (
	"log"
	"net"
	"peer-to-peer/app/shared"
	"sync"
)

/**
 * UDPServerListener is a struct that listens for incoming UDP messages on a specified port.
 * It maintains a map of transport addresses to channels for handling messages.
 * This represent the unique serveur locally, the channels are read by the udp_transport_channel according to the address.
 */
type UDPServerListener struct {
	addr net.UDPAddr
	channels map[TransportAddressKey]ITransportChannel
}

var instance *UDPServerListener
var once sync.Once

func GetUDPServerListener() *UDPServerListener {
	addr := net.UDPAddr{
		Port: shared.UDPPort,
		IP:   net.IPv4zero,
	}
    once.Do(func() {
        instance = &UDPServerListener{
			addr,
			make(map[TransportAddressKey]ITransportChannel),
		}
    })
    return instance
}

func (u *UDPServerListener) Listen() error {
  	//create upd serveur listener on the specified port
	conn, err := net.ListenUDP("udp", &u.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("Serveur de découverte UDP en écoute sur le port", u.addr.Port)

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Erreur de lecture UDP : %v\n", err)
			continue
		}
		message := buffer[:n]
		address := TransportAddress{
			ip:   remoteAddr.IP,
			port: u.addr.Port, // Use the port of the server listener cause the remote as a UDP server listener
		}
		key := address.GetKey()

		if _, exists := u.channels[key]; !exists {
			channel := NewUDPTransportChannel(address)
			RegisterTransportChannel(channel) 
			u.channels[key] = channel
		}

		transportMessage := NewTransportMessage(
			uint32(len(message)),
			message,
			address,
		);

		channel, exist := u.channels[key]
		if !exist {
			log.Printf("No transport channel found for address: %s\n", address.ip.String())
			continue
		}

		err = channel.CollectMessage(transportMessage)
		if err != nil {
			log.Printf("Error collecting message from channel: %v\n", err)
			continue
		}
			
	}
}

