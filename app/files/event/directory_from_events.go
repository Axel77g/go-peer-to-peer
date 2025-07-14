package file_event

import "peer-to-peer/app/shared"

type ShadowFile struct {
	fileName string
	filePath string
	checksum string
}

func (s *ShadowFile) GetName() string {
	return s.fileName
}

func (s *ShadowFile) GetPath() string {
	return s.filePath
}

func (s *ShadowFile) GetChecksum() string {
	return s.checksum
}

func NewShadowFileFromEvent(event shared.FileEvent) *ShadowFile {
	return &ShadowFile{
		fileName: event.FileName,
		filePath: event.FilePath,
		checksum: event.FileChecksum,
	}
}

// applyEventToDirectory applies a file event to the given ShadowDirectory.
// It updates the directory's state based on the type of event:
// - For CreateEvent, it adds a new ShadowFile to the directory.
// - For UpdateEvent, it updates the checksum of the existing ShadowFile.
// - For DeleteEvent, it removes the file from the directory.
//
// Parameters:
//
//	directory - pointer to the ShadowDirectory to be modified.
//	event     - the FileEvent describing the change to apply.
func applyEventToDirectory(directory shared.IDirectory, event shared.FileEvent) {
	switch event.EventType {
	case shared.CreateEvent:
		directory.AddFile(NewShadowFileFromEvent(event))
	case shared.UpdateEvent:
		if file, exists := directory.GetFile(event.FileName); exists {
			file.(*ShadowFile).checksum = event.FileChecksum
		} else {
			//Should not happen
			println("Error: File not found in directory for update event:", event.FileName)
			directory.AddFile(NewShadowFileFromEvent(event))
		}
	case shared.DeleteEvent:
		directory.RemoveFile(event.FileName)
	}
}

func NewShadowDirectory() shared.IDirectory {
	return &shared.Directory{
		DirectoryPath: ":memory:",
		Files:         make(map[string]shared.IFile),
	}
}

func BuildDirectoryFromEvent(events IFileEventCollection) shared.IDirectory {
	directory := NewShadowDirectory()
	iterator := events.GetAll("building directory from events")
	defer iterator.Close()

	for iterator.Next() {
		event, err := iterator.Current()
		if err != nil {
			println("Error reading event:", err)
			return nil
		}
		applyEventToDirectory(directory, event)
	}

	return directory
}
