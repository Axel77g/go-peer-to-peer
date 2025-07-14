package file_watcher

import (
	"peer-to-peer/app/shared"
)

func CompareDirectories(a, b shared.IDirectory) []shared.FileEvent {
	events := []shared.FileEvent{}
	for fileNameA, fileA := range a.GetFiles() {
		if b.HasFile(fileNameA) {
			fileB, _ := b.GetFile(fileNameA)
			if fileB.GetChecksum() != fileA.GetChecksum() { //file updated event
				events = append(events, shared.NewUpdatedFileSystemEvent(fileA))
				continue
			}
		}
		if !b.HasFile(fileNameA) { // file deleted event
			events = append(events, shared.NewDeletedFileSystemEvent(fileA))
			continue
		}
	}

	for fileNameB, fileB := range b.GetFiles() {
		if !a.HasFile(fileNameB) { // file created event
			events = append(events, shared.NewCreateFileSystemEvent(fileB))
		}
	}

	return events
}
