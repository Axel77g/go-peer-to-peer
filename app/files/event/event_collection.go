package file_event

import "peer-to-peer/app/shared"

type IFileEventCollection interface {
	Append(event shared.FileEvent)
	GetAll(reason string) IFileEventIterator
	Merge(collectionB IFileEventCollection) IFileEventCollection
	GetBytesSize() int64
	FromBytes(bytes []byte) error
	SaveToFile(filePath string) error
	GetChecksum() (string, error)
	Debug()
}
