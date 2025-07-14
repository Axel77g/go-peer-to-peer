package file_event

type IFileEventCollection interface {
	Append(event FileEvent)
	GetAll(reason string) IFileEventIterator
	Merge(collectionB IFileEventCollection) IFileEventCollection
	GetBytesSize() int64
	FromBytes(bytes []byte) error
	SaveToFile(filePath string) error
	Debug()
}
