package shared

type Iterator interface {
	Next() bool
	Current() (any, error)
	Reset() error
	Go(int) error
	Size() int
	Close() error
}