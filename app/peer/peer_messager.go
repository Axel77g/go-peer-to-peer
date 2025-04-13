package peer

type PeerMessager interface {
	SendDirectory() error
}