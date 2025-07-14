package file_event

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"peer-to-peer/app/shared"
	"sync/atomic"
)

type JSONLFileEventCollection struct {
	FilePath        string
	activeIterators atomic.Int32
}

func NewJSONLFileEventCollection(filePath string, delete bool) *JSONLFileEventCollection {
	if delete {
		if _, err := os.Stat(filePath); err == nil {
			os.Remove(filePath)
		}
	}

	return &JSONLFileEventCollection{
		FilePath: filePath,
	}
}

func (c *JSONLFileEventCollection) FromBytes(bytes []byte) error {
	file, err := os.OpenFile(c.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

func (c *JSONLFileEventCollection) Append(event shared.FileEvent) {
	if c.activeIterators.Load() > 0 {
		println("Error: Cannot append to collection while iterators are active")
		return
	}

	file, err := os.OpenFile(c.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(event); err != nil {
		println("Error writing event:", err)
	}
}

func (c *JSONLFileEventCollection) GetAll(reason string) IFileEventIterator {
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
	iteratorA := c.GetAll("merging events")
	iteratorB := collectionB.GetAll("merging events")
	defer iteratorA.Close()
	defer iteratorB.Close()

	// On crée une nouvelle collection temporaire
	mergedPath := c.FilePath + "_merged.jsonl"
	mergedCollection := NewJSONLFileEventCollection(mergedPath, true)

	// Pour éviter les doublons
	seen := make(map[string]struct{})

	var (
		eventA shared.FileEvent
		eventB shared.FileEvent
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

func (c *JSONLFileEventCollection) Debug() {
	println("JSONLFileEventCollection Debug:")
	println("FilePath:", c.FilePath)

	fileInfo, err := os.Stat(c.FilePath)
	if err != nil {
		println("Error getting file info:", err)
		return
	}
	println("File Size:", fileInfo.Size(), "bytes")
	println("Active Iterators:", c.activeIterators.Load())

	iterator := c.GetAll("debugging")
	defer iterator.Close()
	if iterator == nil {
		println("No events found in collection.")
		return
	}

	for iterator.Next() {
		event, err := iterator.Current()
		if err != nil {
			println("Error reading event:", err.Error())
			return
		}
		println("Event:", event.Hash)
	}

	println("End of JSONLFileEventCollection Debug\n")

}

func (c *JSONLFileEventCollection) SaveToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		println("Error opening file:", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	iterator := c.GetAll("Saving to file")
	defer iterator.Close()
	defer func() {
		if err := os.Remove(c.FilePath); err != nil {
			println("Error removing old file:", err)
		}
	}()

	for iterator.Next() {
		event, err := iterator.Current()
		if err != nil {
			println("Error reading event:", err)
			return err
		}
		if err := encoder.Encode(event); err != nil {
			println("Error writing event:", err)
			return err
		}
	}

	return nil
}

func (c *JSONLFileEventCollection) GetChecksum() (string, error) {
	if c.activeIterators.Load() > 0 {
		println("Error: Cannot get checksum while iterators are active")
		return "", nil
	}
	c.activeIterators.Add(1)
	defer c.activeIterators.Add(-1)
	file, err := os.Open(c.FilePath)
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
