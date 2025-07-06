package file_event

import "peer-to-peer/app/shared"

// FileEventIteratorAdapter adapts a IFileEventIterator to implement shared.Iterator
type FileEventIteratorAdapter struct {
	iterator IFileEventIterator
}

// NewFileEventIteratorAdapter creates a new adapter
func NewFileEventIteratorAdapter(iterator IFileEventIterator) shared.Iterator {
	return &FileEventIteratorAdapter{
		iterator: iterator,
	}
}

func (f *FileEventIteratorAdapter) Next() bool {
	return f.iterator.Next()
}

func (f *FileEventIteratorAdapter) Current() (any, error) {
	event, err := f.iterator.Current()
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (f *FileEventIteratorAdapter) Reset() error {
	return f.iterator.Reset()
}

func (f *FileEventIteratorAdapter) Go(index int) error {
	return f.iterator.Go(index)
}

func (f *FileEventIteratorAdapter) Size() int {
	return f.iterator.Size()
}

func (f *FileEventIteratorAdapter) Close() error {
	return f.iterator.Close()
}
