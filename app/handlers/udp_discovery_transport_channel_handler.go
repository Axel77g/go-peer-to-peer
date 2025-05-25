package handlers

import (
	"fmt"
	"net"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
	"strconv"
	"strings"
)

type UDPDiscoveryTransportChannel struct {}

func (u *UDPDiscoveryTransportChannel) OnClose(channel peer_comunication.ITransportChannel) {
	peer_comunication.UnregisterTransportChannel(channel)
}

func (u *UDPDiscoveryTransportChannel) OnOpen(channel peer_comunication.ITransportChannel) {
	//peer_comunication.RegisterTransportChannel(channel)
}

func (u *UDPDiscoveryTransportChannel) OnMessage(channel peer_comunication.ITransportChannel, message peer_comunication.TransportMessage) error {
	messageContent := strings.TrimSpace(string(message.GetContent()))
	parts := strings.Split(messageContent, ":")

	if len(parts) > 1 {
			if parts[1] == fmt.Sprintf("%d", shared.SOCKET_ID){
				/* log.Printf(
					"Message from the same socket ID (%d) ignored: %s\n", shared.SOCKET_ID, messageContent) */
				return nil // Ignore the message if it is from the same socket ID
			}
			if parts[0] == "DISCOVER_PEER_REQUEST" {
				address := channel.GetAddress()
				
				peer, existsPeer := peer_comunication.GetPeerByAddress(address)
				if !existsPeer {
					peer_comunication.RegisterTransportChannel(channel)
				    createTCPConnectionForIP(address.GetIP())
				}else{
					_, existsUDP := peer.GetTransportsChannels().GetByAddress(address)
					if !existsUDP {
						peer_comunication.RegisterTransportChannel(channel)
					}

					_, exits := peer.GetTransportsChannels().GetByType("tcp")
					//si aucune connection TCP n'existe pas, on en crée une a la reception d'une discovery request
					if !exits {
						createTCPConnectionForIP(peer.GetIP())
					}
				}
				
				
			}
		}

	return nil
}

func createTCPConnectionForIP(ip net.IP) {
	address := net.JoinHostPort(ip.String(), strconv.Itoa(shared.TCPPort))
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("Failed to create TCP connection for IP %s: %v", ip, err))
	}
	peer_comunication.NewTCPTransportChannel(conn, &TCPTransportChannelHandler{})
}