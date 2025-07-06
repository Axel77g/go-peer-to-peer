package peer_comunication

import "peer-to-peer/app/shared"

type ITransportChannel interface {
	GetProtocol() string
	GetPort() int
	GetAddress() TransportAddress
	Send(content []byte) error
	SendIterator(size uint32, message []byte, iterator shared.Iterator) error
	CollectMessage(TransportMessage) error
	Close() error
	IsAlive() bool
}

type ITransportChannelHandler interface {
	OnOpen(channel ITransportChannel)
	OnClose(channel ITransportChannel)
	OnMessage(channel ITransportChannel, message TransportMessage) error
}
