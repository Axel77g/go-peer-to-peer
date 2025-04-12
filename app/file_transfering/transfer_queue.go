package filetransfering

import (
	"sync"
	"time"
)

type TransferQueue struct {
	mu    sync.Mutex
	Queue []Transfer
}

func NewTransferQueue() *TransferQueue {
	return &TransferQueue{
		Queue: make([]Transfer, 0),
	}
}

func (queue *TransferQueue) Has(transfer Transfer) bool {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	for _, t := range queue.Queue {
		if t.ID == transfer.ID {
			return true
		}
	}
	return false
}

func (queue *TransferQueue) AddTransfer(transfer Transfer) {
	if queue.Has(transfer) {
		return
	}
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.Queue = append(queue.Queue, transfer)
}

func (queue *TransferQueue) Loop() {
	const maxConcurrent = 2
	sem := make(chan int, maxConcurrent)

	for {
		queue.mu.Lock()
		if len(queue.Queue) == 0 {
			queue.mu.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		transfer := queue.Queue[0]
		queue.Queue = queue.Queue[1:]
		queue.mu.Unlock()

		sem <- 1

		go func(t Transfer) {
			defer func() { <-sem }()

			result := t.Start()
			if result.Result {
				println("Transfer réussi:", t.ID)
			} else {
				println("Transfer échoué:", t.ID)
			}
		}(transfer)
	}
}
