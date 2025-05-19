package peer_comunication

type TCPTrasnportMessage struct {
	size    uint32
	content []byte
}

func NewTCPMessage(size uint32, content []byte, isCompressed bool) *TCPTrasnportMessage {
	return &TCPTrasnportMessage{
		size:    size,
		content: content,
	}
}

func (t *TCPTrasnportMessage) getSize() uint32 {
	return t.size
}

func (t *TCPTrasnportMessage) getContent() []byte {
	return t.content
}
