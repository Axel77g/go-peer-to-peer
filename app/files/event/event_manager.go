package file_event

import (
	"fmt"
	"log"
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

func (m *EventManager) AppendEvent(event FileEvent) {
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
func (m *EventManager) MergeAndSave(remoteCollection IFileEventCollection) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	merged := remoteCollection.Merge(m.collection)
	if merged == nil {
		return fmt.Errorf("error merging collections")
	}

	// Save the merged collection to the file
	err := merged.SaveToFile("events.jsonl")
	if err != nil {
		log.Printf("Error saving merged collection: %v", err)
		return fmt.Errorf("failed to save merged collection: %w", err)
	}
	m.collection = NewJSONLFileEventCollection("events.jsonl", false)
	log.Println("Merged and saved events successfully.")
	return nil
}

// SaveCollection explicitly saves the current state of the collection to file.
// This can be used for periodic saves or on application shutdown.
func (m *EventManager) SaveCollection() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.collection.SaveToFile("events.jsonl")
}
