package filesystemwatcher

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
)

/*
*  File Stat
 */

type FileStat struct {
	Name string
	Path string
	Checksum string
}

func (a *FileStat) Compare(b *FileStat) bool {
	if a.Name != b.Name {
		return false
	}

	if a.Checksum != b.Checksum {
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
			if fileB.Checksum != fileA.Checksum { //file updated event
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

func GetFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hasher := md5.New()
    if _, err := io.Copy(hasher, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hasher.Sum(nil)), nil
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


		checksum, err := GetFileChecksum(path)

		if err != nil {
			log.Println("Cannot get the file checksum")
			panic(err)
		}

		fileStat := FileStat{
			Name: dirEntry.Name(),
			Path: path,
			Checksum: checksum,
		}

		directoryList[name] = fileStat
	}

	return DirStat{
		directoryPath,
		directoryList,
	}

}