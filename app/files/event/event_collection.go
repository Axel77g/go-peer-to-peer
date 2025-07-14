package file_event

import "peer-to-peer/app/shared"

type IFileEventCollection interface {
	Append(event shared.FileEvent)
	GetAll(reason string) IFileEventIterator
	Merge(collectionB IFileEventCollection) IFileEventCollection
	GetBytesSize() int64
	FromBytes(bytes []byte) error
	SaveToFile(filePath string) error
	GetChecksum() (string, error)
	Debug()
}

func MergeCollection(collectionA, collectionB, resultCollection IFileEventCollection) IFileEventCollection {
	iteratorA := collectionA.GetAll("merging events")
	iteratorB := collectionB.GetAll("merging events")
	defer iteratorA.Close()
	defer iteratorB.Close()

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
				resultCollection.Append(eventA)
				seen[eventA.Hash] = struct{}{}
			}
			hasA = iteratorA.Next()
		} else if hasB && (!hasA || eventB.Timestamp < eventA.Timestamp) {
			if _, exists := seen[eventB.Hash]; !exists {
				resultCollection.Append(eventB)
				seen[eventB.Hash] = struct{}{}
			}
			hasB = iteratorB.Next()
		} else if hasA && hasB && eventA.Timestamp == eventB.Timestamp {
			// Deux événements avec le même timestamp
			if _, exists := seen[eventA.Hash]; !exists {
				resultCollection.Append(eventA)
				seen[eventA.Hash] = struct{}{}
			}
			if eventA.Hash != eventB.Hash {
				if _, exists := seen[eventB.Hash]; !exists {
					resultCollection.Append(eventB)
					seen[eventB.Hash] = struct{}{}
				}
			}
			hasA = iteratorA.Next()
			hasB = iteratorB.Next()
		}
	}

	return resultCollection
}
