/*
*

---------ATTENTION----------
Code non utilisé pour le moment

*
*/
package filetransfering

import (
	"net"
	"strconv"
)

type TransferResult struct {
	Transfer Transfer
	Result   bool
}

type Transfer struct {
	ID   uint16
	From net.TCPAddr
	To   net.TCPAddr
	File File
}

func NewTransfer(id uint16, from net.TCPAddr, to net.TCPAddr, file File) *Transfer {
	return &Transfer{
		ID:   id,
		From: from,
		To:   to,
		File: file,
	}
}

func (transfer *Transfer) GetStartPayload() []byte {
	return []byte("START:" + strconv.Itoa(int(transfer.ID)) + "," + strconv.Itoa(int(transfer.File.ID)) + "," + transfer.File.Name + "," + strconv.Itoa(int(transfer.File.Size)) + "," + strconv.Itoa(int(transfer.File.ChunkSize)))
}

func (transfer *Transfer) Start() TransferResult {
	//create a TCP connection to the peer
	conn, err := net.DialTCP("tcp", &transfer.From, &transfer.To)
	if err != nil {
		return TransferResult{
			Transfer: *transfer,
			Result:   false,
		}
	}

	defer conn.Close()
	//send the file metadata to the peer
	_, err = conn.Write(transfer.GetStartPayload())
	if err != nil {
		return TransferResult{
			Transfer: *transfer,
			Result:   false,
		}
	}

	//send the file chunks to the peer
	for _, chunk := range transfer.File.Chunks {
		_, err = conn.Write(chunk.GetChunckPayload())
		if err != nil {
			return TransferResult{
				Transfer: *transfer,
				Result:   false,
			}
		}
	}
	return TransferResult{
		Transfer: *transfer,
		Result:   true,
	}
}
