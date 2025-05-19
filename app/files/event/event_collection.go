package file_event

type IFileEventCollection interface {
	Append(event FileEvent)
	GetAll() IFileEventIterator
	Merge(collectionB IFileEventCollection) IFileEventCollection
}