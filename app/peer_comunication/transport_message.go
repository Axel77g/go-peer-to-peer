package peer_comunication

import "time"


type TransportMessage struct {
	size    uint32
	content []byte
	from TransportAddress
	time int64
}

func NewTransportMessage(size uint32, content []byte, from TransportAddress) TransportMessage {
	return TransportMessage{
		size:    size,
		content: content,
		from:   from,
		time:   time.Now().UnixMilli(),
	}
}

func (t *TransportMessage) GetSize() uint32 {
	return t.size
}

func (t *TransportMessage) GetContent() []byte {
	return t.content
}

func (t *TransportMessage) GetFrom() TransportAddress {
	return t.from
}

func (t *TransportMessage) GetTime() int64 {
	return t.time
}