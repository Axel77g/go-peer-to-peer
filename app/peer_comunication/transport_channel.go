package peer_comunication

type ITransportChannel interface {
	GetPort() int
	Send(content []byte) error
	Read() (ITransportMessage, error)
	Close() error
}