package file_event
type IFileEventIterator interface {
	Next() bool
	Current() (FileEvent, error)
	Reset() error
	Go(int) error
	Size() int
	Close() error
}