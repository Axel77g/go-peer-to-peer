package peer_comunication

type UDPTransportMessage struct {
	content []byte
	size    uint32
	from    TransportAddress
}

func NewUDPTransportMessage(content []byte, from TransportAddress) *UDPTransportMessage {
	if len(content) > 0xFFFFFFFF {
		panic("Message size exceeds maximum limit")
	}
	return &UDPTransportMessage{
		content: content,
		size:    uint32(len(content)),
		from:    from,
	}
}

func (m *UDPTransportMessage) getSize() uint32 {
	return m.size
}

func (m *UDPTransportMessage) getContent() []byte {
	return m.content
}

func (m *UDPTransportMessage) getFrom() TransportAddress {
	return m.from
}
