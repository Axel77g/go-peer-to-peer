package peer_comunication

type ITransportMessage interface {
	getSize() uint32
	getContent() []byte
}