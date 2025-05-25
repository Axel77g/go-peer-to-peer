package file_event

import (
	"bytes"
	"encoding/gob"
	"os"
)

type ByteEventIterator struct {
	src         IFileEventIterator
	currentData []byte
}

func NewByteEventIterator(src IFileEventIterator) *ByteEventIterator {
	return &ByteEventIterator{
		src:        src,
	}
}

func serialize(event FileEvent) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(event)
	return buf.Bytes(), err
}

func (b *ByteEventIterator) Next() bool {
	if b.src.Next() {
		currentEvent, err := b.src.Current()
		if err != nil {
			return false
		}
		data, err := serialize(currentEvent)
		if err != nil {
			return false
		}
		b.currentData = data
		return true
	}
	return false
}
func (b *ByteEventIterator) Current() ([]byte, error) {
	if len(b.currentData) == 0 {
		return nil, os.ErrInvalid
	}
	return b.currentData, nil
}
func (b *ByteEventIterator) Reset() error {
	if err := b.src.Reset(); err != nil {
		return err
	}
	b.currentData = nil
	return nil
}
func (b *ByteEventIterator) Go(n int) error {
	if err := b.src.Go(n); err != nil {
		return err
	}
	if !b.Next() {
		return os.ErrInvalid
	}
	return nil
}
func (b *ByteEventIterator) Size() int {
	return b.src.Size()
}
func (b *ByteEventIterator) Close() error {
	if err := b.src.Close(); err != nil {
		return err
	}
	b.currentData = nil
	return nil
}
// OnIteratorClose is called when the iterator is closed
func (b *ByteEventIterator) OnIteratorClose() {
	if closer, ok := b.src.(IFileEventIteratorCloser); ok {
		closer.OnIteratorClose()
	}
}
// GetBytesSize returns the size of the current data in bytes
func (b *ByteEventIterator) GetBytesSize() int64 {
	if len(b.currentData) == 0 {
		return 0
	}
	return int64(len(b.currentData))
}
