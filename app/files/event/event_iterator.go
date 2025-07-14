package file_event

import "peer-to-peer/app/shared"

type IFileEventIterator interface {
	Next() bool
	Current() (shared.FileEvent, error)
	Reset() error
	Go(int) error
	Size() int
	Close() error
}
