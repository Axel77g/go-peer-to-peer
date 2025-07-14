package file_event

import (
	"fmt"
	"log"
	"peer-to-peer/app/peer_comunication"
	"peer-to-peer/app/shared"
	"sync"
)

var (
	eventManagerInstance *EventManager
	once                 sync.Once
)

// EventManager manages concurrent access to the main event collection file.
type EventManager struct {
	collection IFileEventCollection
	mu         sync.Mutex
}

func GetEventManager() *EventManager {
	once.Do(func() {
		// Initialize the collection. Load existing events if the file exists.
		// The second argument 'false' means it will attempt to load from the file.
		collection := NewJSONLFileEventCollection("events.jsonl", false)

		eventManagerInstance = &EventManager{
			collection: collection,
		}
		log.Println("Event manager initialized.")
	})
	return eventManagerInstance
}

func (m *EventManager) AppendEvent(event shared.FileEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collection.Append(event)
	// Consider adding logic here to mark the collection as dirty
	// and trigger a save periodically or on a separate goroutine.
}

func (m *EventManager) Lock() {
	m.mu.Lock()
}

func (m *EventManager) Unlock() {
	m.mu.Unlock()
}

func (m *EventManager) GetCollection() IFileEventCollection {
	return m.collection
}

// MergeAndSave merges a remote collection into the managed collection and saves the result to file.
// This operation is thread-safe.
func (m *EventManager) MergeAndSave(remoteCollection IFileEventCollection) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	baseChecksum, err := m.collection.GetChecksum()
	if err != nil {
		log.Printf("Error getting base collection checksum: %v", err)
		return false, fmt.Errorf("failed to get base collection checksum: %w", err)
	}

	remoteChecksum, err := remoteCollection.GetChecksum()
	if err != nil {
		log.Printf("Error getting remote collection checksum: %v", err)
		return false, fmt.Errorf("failed to get remote collection checksum: %w", err)
	}

	if baseChecksum == remoteChecksum {
		log.Println("No changes detected, skipping merge.")
		return false, nil
	}

	merged := remoteCollection.Merge(m.collection)
	if merged == nil {
		return false, fmt.Errorf("error merging collections")
	}

	// Save the merged collection to the file
	err = merged.SaveToFile("events.jsonl")
	if err != nil {
		log.Printf("Error saving merged collection: %v", err)
		return false, fmt.Errorf("failed to save merged collection: %w", err)
	}
	m.collection = NewJSONLFileEventCollection("events.jsonl", false)
	log.Println("Merged and saved events successfully.")
	return true, nil
}

// SaveCollection explicitly saves the current state of the collection to file.
// This can be used for periodic saves or on application shutdown.
func (m *EventManager) SaveCollection() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.collection.SaveToFile("events.jsonl")
}

func (m *EventManager) BroadcastEvents() {
	log.Println("Broadcasting events to all peers.")
	m.mu.Lock()
	defer m.mu.Unlock()
	iterator := m.collection.GetAll("broadcasting events")
	defer iterator.Close()
	peer_comunication.BroadcastIterator([]byte("PUSH_EVENTS"), NewFileEventIteratorAdapter(iterator), "tcp")
}
