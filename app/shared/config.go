package shared

import "math/rand"

var (
	TCPPort = 9998
	UDPPort = 9999
	SHARED_DIRECTORY = "./Shared"
	SOCKET_ID = rand.Intn(10000000)
)