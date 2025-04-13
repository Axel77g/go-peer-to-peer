package shared

import (
	"net"
	"sync"
)
type Socket interface {
	ListenMessage(conn net.Conn, mu *sync.Mutex)
	Send(message any) (bool, error)
	GetConn() (net.Conn, error)
}