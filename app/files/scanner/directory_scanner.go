package file_scanner

import (
	"log"
	"os"
	"path/filepath"
	"peer-to-peer/app/shared"
)

func Scan(directoryPath string) shared.IDirectory {
	directoryList := make(map[string]shared.IFile)

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

		fileStat, err := shared.NewFile(
			dirEntry.Name(),
			path,
		)

		if err != nil {
			log.Println("Cannot create file stat")
			panic(err)
		}

		directoryList[name] = fileStat
	}

	return shared.NewDirectory(directoryPath, directoryList)
}