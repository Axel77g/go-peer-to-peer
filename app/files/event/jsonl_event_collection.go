package file_event

import (
	"encoding/json"
	"os"
	"sync/atomic"
)

type JSONLFileEventCollection struct {
	FilePath        string
	FileMode        int
	activeIterators atomic.Int32
}

func NewJSONLFileEventCollection(filePath string, mode ...int) *JSONLFileEventCollection {
	fileMode := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	if len(mode) > 0 {
		if len(mode) == 1 {
			fileMode = mode[0]
		} else {
			fileMode = 0
			for _, m := range mode {
				fileMode |= m
			}
		}
	}
	return &JSONLFileEventCollection{
		FilePath: filePath,
		FileMode: fileMode,
	}
}

func(c *JSONLFileEventCollection) FromBytes(bytes []byte) error {
	file, err := os.OpenFile(c.FilePath, c.FileMode, 0644)
	if err != nil {
		println("Error opening file:", err)
	}
	defer file.Close()
	
	_, err = file.Write(bytes)
	if err != nil {
		println("Error writing bytes to file:", err)
	}

	return nil
}

func (c *JSONLFileEventCollection) OnIteratorClose() {
	current := c.activeIterators.Add(-1)
	if current < 0 {
		c.activeIterators.Store(0)
	}
}

func (c *JSONLFileEventCollection) Append(event FileEvent) {
	if c.activeIterators.Load() > 0 {
		println("Error: Cannot append to collection while iterators are active")
		return
	}

	file, err := os.OpenFile(c.FilePath, c.FileMode, 0644)
	if err != nil {
		println("Error opening file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(event); err != nil {
		println("Error writing event:", err)
	}
}

func (c *JSONLFileEventCollection) GetAll() IFileEventIterator {
	c.activeIterators.Add(1)
	iterator, err := NewJSONLFileEventIterator(c.FilePath, c)
	if err != nil {
		println("Error: ", err)
		return nil
	}
	return iterator
}

func (c *JSONLFileEventCollection) Merge(collectionB IFileEventCollection) IFileEventCollection {
	//redigier ici la merge logic
	iteratorA := c.GetAll()
	iteratorB := collectionB.GetAll()
	defer iteratorA.Close()
	defer iteratorB.Close()

	// On crée une nouvelle collection temporaire
	mergedPath := c.FilePath + "_merged.jsonl"
	mergedCollection := NewJSONLFileEventCollection(mergedPath,os.O_WRONLY, os.O_CREATE)

	// Pour éviter les doublons
	seen := make(map[string]struct{})

	var (
		eventA FileEvent
		eventB FileEvent
		hasA   = iteratorA.Next()
		hasB   = iteratorB.Next()
	)

	for hasA || hasB {
		var err error

		if hasA {
			eventA, err = iteratorA.Current()
			if err != nil {
				println("Error reading eventA:", err)
				return nil
			}
		}
		if hasB {
			eventB, err = iteratorB.Current()
			if err != nil {
				println("Error reading eventB:", err)
				return nil
			}
		}

		if hasA && (!hasB || eventA.Timestamp < eventB.Timestamp) {
			if _, exists := seen[eventA.Hash]; !exists {
				mergedCollection.Append(eventA)
				seen[eventA.Hash] = struct{}{}
			}
			hasA = iteratorA.Next()
		} else if hasB && (!hasA || eventB.Timestamp < eventA.Timestamp) {
			if _, exists := seen[eventB.Hash]; !exists {
				mergedCollection.Append(eventB)
				seen[eventB.Hash] = struct{}{}
			}
			hasB = iteratorB.Next()
		} else if hasA && hasB && eventA.Timestamp == eventB.Timestamp {
			// Deux événements avec le même timestamp
			if _, exists := seen[eventA.Hash]; !exists {
				mergedCollection.Append(eventA)
				seen[eventA.Hash] = struct{}{}
			}
			if eventA.Hash != eventB.Hash {
				if _, exists := seen[eventB.Hash]; !exists {
					mergedCollection.Append(eventB)
					seen[eventB.Hash] = struct{}{}
				}
			}
			hasA = iteratorA.Next()
			hasB = iteratorB.Next()
		}
	}

	return mergedCollection
}

func (c *JSONLFileEventCollection) GetBytesSize() int64 {
	fileInfo, err := os.Stat(c.FilePath)
	if err != nil {
		println("Error getting file size:", err)
		return 0
	}
	return fileInfo.Size()
}