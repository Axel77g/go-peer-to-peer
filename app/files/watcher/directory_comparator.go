package file_watcher

import (
	file_event "peer-to-peer/app/files/event"
	"peer-to-peer/app/shared"
)

func CompareDirectories(a, b shared.IDirectory) []file_event.FileEvent {
	events := []file_event.FileEvent{}
	for fileNameA, fileA := range a.GetFiles() {
		if b.HasFile(fileNameA) {
			fileB, _ := b.GetFile(fileNameA)
			if fileB.GetChecksum() != fileA.GetChecksum() { //file updated event
				events = append(events, file_event.NewUpdatedFileSystemEvent(fileA))
				continue
			}
		}
		if !b.HasFile(fileNameA) { // file deleted event
			events = append(events, file_event.NewDeletedFileSystemEvent(fileA))
			continue
		}
	}

	for fileNameB, fileB := range b.GetFiles() {
		if !a.HasFile(fileNameB) { // file created event
			events = append(events, file_event.NewCreateFileSystemEvent(fileB))
		}
	}

	return events
}
