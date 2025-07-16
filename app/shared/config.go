package shared

import "math/rand"

var (
	TCPPort               int    = 9998
	UDPPort               int    = 9999
	SHARED_DIRECTORY      string = "./shared"
	EVENT_COLLECTION_FILE string = "./events.jsonl"
	SOCKET_ID             int    = rand.Intn(10000000)
)
