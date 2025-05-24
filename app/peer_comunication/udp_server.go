package peer_comunication

import (
	"errors"
	"log"
	"net"
)

type UDPServer struct {
	port                      int
	transportsChannelsMessage map[*net.IP]chan ITransportMessage
	stop                      chan struct{}
}

func NewUDPServer(port int) UDPServer {

	server := UDPServer{
		port:                      port,
		transportsChannelsMessage: make(map[*net.IP]chan ITransportMessage),
		stop:                      make(chan struct{}),
	}
	return server
}

func (u *UDPServer) readMessage(conn *net.UDPConn) (ITransportMessage, error) {
	messageSize := make([]byte, 4)
	size, remoteAddr, err := conn.ReadFromUDP(messageSize)
	if err != nil {
		log.Println("Erreur lors de la lecture du message UDP :", err)
		return nil, err
	}

	if size != 4 {
		return nil, errors.New("message size is not 4 bytes")
	}

	messageContent := make([]byte, int(messageSize[0]))
	_, _, err = conn.ReadFromUDP(messageContent)
	if err != nil {
		log.Println("Erreur lors de la lecture du message UDP :", err)
		return nil, err
	}

	addr := TransportAddress{
		ip:   remoteAddr.IP,
		port: remoteAddr.Port,
	}
	message := NewUDPTransportMessage(messageContent, addr)
	return message, nil
}

func (u *UDPServer) listen() {
	addr := net.UDPAddr{
		Port: u.port,
		IP:   net.IPv4zero,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Printf("Erreur lors de l'écoute UDP : %v\n", err)
		return
	}

	defer conn.Close()
	log.Println("UDP server listening on port", u.port)

	for {
		select {
		case <-u.stop:
			return
		default:
			message, err := u.readMessage()
			if err != nil {
				continue
			}
			u.collectMessage(message)
		}
	}
}

func (u *UDPServer) collectMessage(message ITransportMessage) {
	addr := message.getFrom()
	channel, exist := u.transportsChannelsMessage[&addr.ip]
	if !exist {
		messageChan := make(chan ITransportMessage, 100)
		u.transportsChannelsMessage[&addr.ip] = messageChan
		channel := u.transportsChannelsMessage[&addr.ip]
		channel <- message
	} else {
		if len(channel) == 100 {
			<-channel
		}
		channel <- message
	}
}
