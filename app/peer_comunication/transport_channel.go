package peer_comunication

type ITransportChannel interface {
	GetPort() int
	GetAddress() TransportAddress
	Send(content []byte) error
	CollectMessage(TransportMessage) error
	Read() (TransportMessage, error)
	Close() error
}
