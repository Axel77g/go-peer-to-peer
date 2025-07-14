package file_event

import (
	"crypto/sha256"
	"encoding/hex"
	"peer-to-peer/app/shared"
	"strconv"
	"time"
)

type EventType int

const (
	CreateEvent EventType = iota
	UpdateEvent
	DeleteEvent
)

type FileEvent struct {
	Hash         string
	Timestamp    uint64
	EventType    EventType
	FileName	string
	FilePath     string
	FileChecksum *string
}

func getCurrentTimestamp() uint64 {
	return uint64(time.Now().UnixNano())
}

func generateHash(filePath string, timestamp uint64, eventType EventType) string {
	data := filePath + strconv.FormatUint(timestamp, 10) + strconv.FormatUint(uint64(eventType), 10)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func NewCreateFileEvent(file shared.IFile, eventType EventType) FileEvent {
	timestamp := getCurrentTimestamp()
	hash := generateHash(file.GetPath(), timestamp, eventType)
	return FileEvent{
		Hash:      hash,
		Timestamp: timestamp,
		EventType: eventType,
		FileName:  file.GetPath(),
		FilePath:  file.GetPath(),
		FileChecksum: file.GetChecksum(),
	}
}

func NewCreateFileSystemEvent(file shared.IFile) FileEvent {
	return NewCreateFileEvent(file, CreateEvent)
}

func NewUpdatedFileSystemEvent(file shared.IFile) FileEvent {
	return NewCreateFileEvent(file, UpdateEvent)
}

func NewDeletedFileSystemEvent(file shared.IFile) FileEvent {
	return NewCreateFileEvent(file, DeleteEvent)
}