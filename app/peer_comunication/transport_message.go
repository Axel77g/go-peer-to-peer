package peer_comunication

type TransportMessage struct {
	size    uint32
	content []byte
	from TransportAddress
}

func NewTransportMessage(size uint32, content []byte, from TransportAddress) TransportMessage {
	return TransportMessage{
		size:    size,
		content: content,
		from:   from,
	}
}

func (t *TransportMessage) getSize() uint32 {
	return t.size
}

func (t *TransportMessage) getContent() []byte {
	return t.content
}

func (t *TransportMessage) getFrom() TransportAddress {
	return t.from
}