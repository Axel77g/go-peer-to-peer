package peer_comunication

type ITransportChannel interface {
	GetProtocol() string
	GetPort() int
	GetAddress() TransportAddress
	Send(content []byte) error
	CollectMessage(TransportMessage) error
	Close() error
	IsAlive() bool
}

type ITransportChannelHandler interface {
	OnOpen(channel ITransportChannel)
	OnClose(channel ITransportChannel)
	OnMessage(channel ITransportChannel, message TransportMessage) error
}
