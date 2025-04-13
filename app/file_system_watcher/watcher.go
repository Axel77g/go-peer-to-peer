package filesystemwatcher

import (
	"log"
	"time"
)

type Watcher struct {
	DirectoryPath string
	Cooldown      time.Duration
	Events        chan FileSystemEvent
}

func NewWatcher(dirPath string, cooldown time.Duration, event chan FileSystemEvent) Watcher {
	return Watcher{
		dirPath,
		cooldown,
		event,
	}
}



// @blocking thread
func (watcher *Watcher) Listen() {
	log.Printf("Listening for file system events in %s\n", watcher.DirectoryPath)
	baseDirStat := GetFileDirStat(watcher.DirectoryPath)
	ticker := time.NewTicker(watcher.Cooldown)
	defer ticker.Stop()
	for range ticker.C {
		newFileDirStat := GetFileDirStat(watcher.DirectoryPath)
		events := baseDirStat.Compare(&newFileDirStat)
		baseDirStat = newFileDirStat
		if len(events) == 0 {
			continue
		}

		log.Printf("File system events detected: %v\n", events)

		for _, event := range events {
			watcher.Events <- event
		}
	}
}
