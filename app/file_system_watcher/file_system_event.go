package filesystemwatcher

type FileSystemEvent struct {
	EventType uint
	FilePath  string
}

func NewCreateFileSystemEvent(filePath string) FileSystemEvent {
	return FileSystemEvent{
		1,
		filePath,
	}
}

func NewUpdatedFileSystemEvent(filePath string) FileSystemEvent {
	return FileSystemEvent{
		2,
		filePath,
	}
}

func NewDeletedFileSystemEvent(filePath string) FileSystemEvent {
	return FileSystemEvent{
		3,
		filePath,
	}
}
