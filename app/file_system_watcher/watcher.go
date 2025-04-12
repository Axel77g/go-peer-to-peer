package filesystemwatcher

import (
	"log"
	"os"
	"path/filepath"
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

func GetFileDirStat(directoryPath string) DirStat {
	directoryList := make(map[string]FileStat)

	directoryReadResult, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Println("Cannot read the directry")
		panic(err)
	}

	for _, dirEntry := range directoryReadResult {
		if dirEntry.IsDir() {
			continue
		}
		name := dirEntry.Name()
		path := filepath.Join(directoryPath, name)

		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Printf("Cannot get file info from the file %s\n", name)
			continue
		}

		fileStat := FileStat{
			Name: dirEntry.Name(),
			Path: path,
			Size: uint64(fileInfo.Size()),
		}

		directoryList[name] = fileStat
	}

	return DirStat{
		directoryPath,
		directoryList,
	}

}

// @blocking thread
func (watcher *Watcher) Listen() {
	log.Printf("Listening for file system events in %s\n", watcher.DirectoryPath)
	baseDirStat := GetFileDirStat(watcher.DirectoryPath)
	ticker := time.NewTicker(watcher.Cooldown)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
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
}
