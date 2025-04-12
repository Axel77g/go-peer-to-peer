package tcpcomunication

import "encoding/json"

var (
	MESSAGE_TYPE_HELLO = "HELLO"
	MESSAGE_TYPE_FILE_DIR = "FILE_DIR"
)

type TCPMessageJSON struct {
	Type string `json:"type"`
	Payload  json.RawMessage `json:"payload"`
}

type HelloMessage struct {
	PeerID string `json:"socket_id"`
}

type FileDirMessage struct {
	Files []string `json:"files"`
}