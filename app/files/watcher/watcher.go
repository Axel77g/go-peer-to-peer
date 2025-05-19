package file_watcher

import (
	"log"
	file_event "peer-to-peer/app/files/event"
	file_scanner "peer-to-peer/app/files/scanner"
	"time"
)

type Watcher struct {
	DirectoryPath string
	Cooldown      time.Duration
	Events        chan file_event.FileEvent
}

func NewWatcher(dirPath string, cooldown time.Duration, event chan file_event.FileEvent) Watcher {
	return Watcher{
		dirPath,
		cooldown,
		event,
	}
}

// @blocking thread
func (watcher *Watcher) Listen() {
	log.Printf("Listening for file system events in %s\n", watcher.DirectoryPath)
	baseDirStat := file_scanner.Scan(watcher.DirectoryPath)
	ticker := time.NewTicker(watcher.Cooldown)
	defer ticker.Stop()
	for range ticker.C {
		newFileDirStat := file_scanner.Scan(watcher.DirectoryPath)
		events := CompareDirectories(baseDirStat, newFileDirStat)
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
