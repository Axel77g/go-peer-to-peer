package file_event

import (
	"os"
	"testing"
	"time"
)

type MockFile struct {
	FilePath string
}

func (m MockFile) GetPath() string {
	return m.FilePath
}

func (m MockFile) GetName() string {
	return m.FilePath
}

func (m MockFile) GetChecksum() *string {
	checksum := "mocked_checksum"
	return &checksum
}


func TestJSONLFileEventCollection_Append_GetAll(t *testing.T) {
	tmpFile := "test_events_A.jsonl"
	defer os.Remove(tmpFile)
	collection := NewJSONLFileEventCollection(tmpFile)

	event1 := NewCreateFileEvent(MockFile{
		FilePath: "file1.txt",
	}, 1)
	event2 := NewCreateFileEvent(MockFile{
		FilePath: "file2.txt",
	}, 2)

	collection.Append(event1)
	collection.Append(event2)

	it := collection.GetAll()
	defer it.Close()

	var events []FileEvent
	for it.Next() {
		e, err := it.Current()
		if err != nil {
			t.Errorf("Erreur lors de la lecture de l'événement : %v", err)
		}
		events = append(events, e)
	}

	if len(events) != 2 {
		t.Errorf("Attendu 2 événements, obtenu %d", len(events))
	}

	if events[0].Hash != event1.Hash || events[1].Hash != event2.Hash {
		t.Errorf("Les événements lus ne correspondent pas aux événements ajoutés")
	}
}

func TestJSONLFileEventCollection_Merge(t *testing.T) {
	tmpFileA := "test_events_A.jsonl"
	tmpFileB := "test_events_B.jsonl"
	defer os.Remove(tmpFileA)
	defer os.Remove(tmpFileB)
	collectionA := NewJSONLFileEventCollection(tmpFileA)
	collectionB := NewJSONLFileEventCollection(tmpFileB)

	// Événement en commun (même hash)
	sharedEvent := NewCreateFileEvent(MockFile{
		FilePath: "shared.txt",
	}, 1)

	// Événements uniques
	eventA := NewCreateFileEvent(MockFile{
		FilePath: "onlyA.txt",
	}, 1)
	eventA2 := NewCreateFileEvent(MockFile{
		FilePath: "onlyA2.txt",
	}, 1)
	time.Sleep(10 * time.Millisecond)
	eventB := NewCreateFileEvent(MockFile{
		FilePath: "onlyB.txt",
	}, 1)

	collectionA.Append(sharedEvent)
	collectionA.Append(eventA)
	collectionA.Append(eventA2)

	collectionB.Append(sharedEvent)
	collectionB.Append(eventB)

	// On merge les deux
	merged := collectionA.Merge(collectionB)
	mergedPath := merged.(*JSONLFileEventCollection).FilePath
	defer os.Remove(mergedPath)

	it := merged.GetAll()
	defer it.Close()

	hashSet := make(map[string]struct{})
	count := 0
	for it.Next() {
		e, err := it.Current()
		if err != nil {
			t.Fatal("Erreur lors de la lecture de l'événement :", err)
		}
		if _, exists := hashSet[e.Hash]; exists {
			t.Errorf("Duplicate hash found: %s", e.Hash)
		}
		hashSet[e.Hash] = struct{}{}
		count++
	}

	if count != 4 {
		t.Errorf("Attendu 3 événements distincts, obtenu %d", count)
	}

	//assert that last merged is the hash of B
	it = merged.GetAll()
	defer it.Close()
	it.Go(it.Size() - 1)
	lastEvent, err := it.Current()
	if err != nil {
		t.Fatal("Erreur lors de la lecture de l'événement :", err)
	}
	if lastEvent.Hash != eventB.Hash {
		t.Errorf("Le dernier événement fusionné n'est pas celui de B, attendu %s, obtenu %s", eventB.Hash, lastEvent.Hash)
	}

}