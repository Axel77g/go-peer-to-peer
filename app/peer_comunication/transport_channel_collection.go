package peer_comunication

import "sync"

type TransportChannelCollection struct {
	channels []ITransportChannel
	mutex             sync.Mutex
}

func (t *TransportChannelCollection) Add(channel ITransportChannel) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.channels = append(t.channels, channel)
}

func (t *TransportChannelCollection) Remove(channel ITransportChannel) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for i, ch := range t.channels {
		addr := ch.GetAddress()
		addr2 := channel.GetAddress()
		if addr.GetKey() == addr2.GetKey() {
			t.channels = append(t.channels[:i], t.channels[i+1:]...)
			break
		}
	}
}

func (t *TransportChannelCollection) GetAll() []ITransportChannel {
	return t.channels
}

func (t *TransportChannelCollection) GetByAddress(address TransportAddress) (ITransportChannel, bool) {
	for _, ch := range t.channels {
		addr := ch.GetAddress()
		if addr.GetKey() == address.GetKey() {
			return ch, true
		}
	}
	return nil, false
}

func (t *TransportChannelCollection) GetByType(channelType string) (ITransportChannel, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, ch := range t.channels {
		if ch.GetProtocol() == channelType {
			return ch, true
		}
	}
	return nil, false

}

func (t *TransportChannelCollection) Size() int {
	return len(t.channels)
}