package filesystemwatcher

/*
*  File Stat
 */

type FileStat struct {
	Name string
	Path string
	Size uint64
}

func (a *FileStat) Compare(b *FileStat) bool {
	if a.Name != b.Name {
		return false
	}

	if a.Size != b.Size {
		return false
	}

	return true
}

/*
*  Directory Stat
 */

type DirStat struct {
	directoryPath string
	Files         map[string]FileStat
}

func (dirStat *DirStat) HasFile(fileName string) bool {
	_, exists := dirStat.Files[fileName]
	return exists
}

func (a *DirStat) Compare(b *DirStat) []FileSystemEvent {
	events := []FileSystemEvent{}

	for fileNameA, fileA := range a.Files {
		if b.HasFile(fileNameA) {
			fileB := b.Files[fileNameA]
			if fileB.Size != fileA.Size { //file updated event
				events = append(events, NewUpdatedFileSystemEvent(fileA.Path))
				continue
			}
		}
		if !b.HasFile(fileNameA) { // file deleted event
			events = append(events, NewDeletedFileSystemEvent(fileA.Path))
			continue
		}
	}

	for fileNameB, fileB := range b.Files {
		if !a.HasFile(fileNameB) { // file created event
			events = append(events, NewCreateFileSystemEvent(fileB.Path))
		}
	}

	return events
}
