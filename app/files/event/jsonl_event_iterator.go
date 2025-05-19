package file_event

import (
	"bufio"
	"encoding/json"
	"os"
)

type JSONLFileEventIterator struct {
	filePath    string
	currentLine int
	size        int
	file        *os.File
	scanner     *bufio.Scanner
	buffer      string
	collection  *JSONLFileEventCollection
}

func NewJSONLFileEventIterator(filePath string, collection *JSONLFileEventCollection) (*JSONLFileEventIterator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		file.Close()
		return nil, err
	}

	// Rewind the file for reading
	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, err
	}

	return &JSONLFileEventIterator{
		filePath:    filePath,
		currentLine: -1,
		size:        lineCount,
		file:        file,
		scanner:     bufio.NewScanner(file),
		collection: collection,
	}, nil
}

func (j *JSONLFileEventIterator) Next() bool {
	if j.currentLine+1 < j.size {
		j.currentLine++
		if j.scanner.Scan() {
			j.buffer = j.scanner.Text()
			return true
		}
	}
	return false
}

func (j *JSONLFileEventIterator) Current() (FileEvent, error) {
	if j.currentLine < 0 || j.currentLine >= j.size {
		return FileEvent{}, os.ErrInvalid
	}

	var event FileEvent
	if err := json.Unmarshal([]byte(j.buffer), &event); err != nil {
		return FileEvent{}, err
	}
	return event, nil
}

func (j *JSONLFileEventIterator) Reset() error {
	if _, err := j.file.Seek(0, 0); err != nil {
		return err
	}
	j.scanner = bufio.NewScanner(j.file)
	j.currentLine = -1
	return nil
}

func (j *JSONLFileEventIterator) Go(index int) error {
	if index < 0 || index >= j.size {
		return os.ErrInvalid
	}
	if err := j.Reset(); err != nil {
		return err
	}
	for i := 0; i <= index; i++ {
		if !j.Next() {
			return os.ErrInvalid
		}
	}
	return nil
}

func (j *JSONLFileEventIterator) Close() error {
	if j.file != nil {
		err := j.file.Close()
		j.file = nil
		return err
	}
	j.collection.OnIteratorClose()
	return nil
}

func (j *JSONLFileEventIterator) Size() int {
	return j.size
}