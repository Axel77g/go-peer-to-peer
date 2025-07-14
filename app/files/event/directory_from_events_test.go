package file_event

import (
	"bytes"
	"errors"
	"io"
	"os"
	"peer-to-peer/app/shared"
	"testing"
)

// MockFileEventIterator is a mock implementation of the IFileEventIterator interface for testing.
type MockFileEventIterator struct {
	events    []shared.FileEvent
	index     int
	closed    bool
	forceErr  bool
	errorStep int
}

func (m *MockFileEventIterator) Next() bool {
	if m.closed || m.index >= len(m.events) {
		return false
	}
	m.index++
	return true
}

func (m *MockFileEventIterator) Current() (shared.FileEvent, error) {
	if m.closed {
		return shared.FileEvent{}, errors.New("iterator is closed")
	}

	if m.forceErr && m.index == m.errorStep {
		return shared.FileEvent{}, errors.New("forced error for testing")
	}

	if m.index <= 0 || m.index > len(m.events) {
		return shared.FileEvent{}, errors.New("invalid iterator position")
	}

	return m.events[m.index-1], nil
}

func (m *MockFileEventIterator) Close() error {
	m.closed = true
	return nil
}

// Additional methods to fully implement the interface pattern based on JSONLFileEventIterator
func (m *MockFileEventIterator) Reset() error {
	if m.closed {
		return errors.New("iterator is closed")
	}
	m.index = -1
	return nil
}

func (m *MockFileEventIterator) Go(index int) error {
	if m.closed {
		return errors.New("iterator is closed")
	}

	if index < 0 || index >= len(m.events) {
		return os.ErrInvalid
	}

	m.index = index
	return nil
}

func (m *MockFileEventIterator) Size() int {
	return len(m.events)
}

// MockFileEventCollection is a mock implementation of the IFileEventCollection interface for testing.
type MockFileEventCollection struct {
	events    []shared.FileEvent
	forceErr  bool
	errorStep int
}

func (m *MockFileEventCollection) GetAll(reason string) IFileEventIterator {
	return &MockFileEventIterator{
		events:    m.events,
		index:     0,
		closed:    false,
		forceErr:  m.forceErr,
		errorStep: m.errorStep,
	}
}

// Implement the remaining methods of IFileEventCollection
func (m *MockFileEventCollection) Append(event shared.FileEvent) {
	m.events = append(m.events, event)
}

func (m *MockFileEventCollection) Merge(collectionB IFileEventCollection) IFileEventCollection {
	result := &MockFileEventCollection{
		events:    make([]shared.FileEvent, len(m.events)),
		forceErr:  m.forceErr,
		errorStep: m.errorStep,
	}

	// Copy events from this collection
	copy(result.events, m.events)

	// Add events from the other collection
	if collectionB != nil {
		iterator := collectionB.GetAll("merging in mock")
		defer iterator.Close()

		for iterator.Next() {
			event, err := iterator.Current()
			if err == nil {
				result.events = append(result.events, event)
			}
		}
	}

	return result
}

func (m *MockFileEventCollection) GetBytesSize() int64 {
	return int64(len(m.events))
}

func (m *MockFileEventCollection) FromBytes(data []byte) error {
	return nil // Not used by directory_from_events.go
}

func (m *MockFileEventCollection) SaveToFile(filePath string) error {
	return nil // Not used by directory_from_events.go
}

func (m *MockFileEventCollection) GetChecksum() (string, error) {
	return "", nil // Not used by directory_from_events.go
}

func (m *MockFileEventCollection) Debug() {
	// No-op for testing
}

type bytesScanner struct {
	reader    *bytes.Reader
	data      []byte
	buf       []byte
	start     int64
	end       int64
	err       error
	bytesRead bool
	done      bool
}

func (s *bytesScanner) Scan() bool {
	if s.done {
		return false
	}

	line, err := s.readLine()
	if err != nil {
		if err != io.EOF {
			s.err = err
		}
		s.done = true
		return false
	}

	s.data = line
	return true
}

func (s *bytesScanner) readLine() ([]byte, error) {
	var line []byte

	for {
		n, err := s.reader.Read(s.buf)
		if n == 0 && err != nil {
			return line, err
		}

		for i := 0; i < n; i++ {
			if s.buf[i] == '\n' {
				line = append(line, s.buf[:i]...)
				s.reader.Seek(int64(-(n-i-1)), io.SeekCurrent)
				return line, nil
			}
		}

		line = append(line, s.buf[:n]...)
	}
}

func (s *bytesScanner) Bytes() []byte {
	return s.data
}

func (s *bytesScanner) Err() error {
	return s.err
}

// OnIteratorClose is a mock implementation to match JSONLFileEventCollection's method
func (m *MockFileEventCollection) OnIteratorClose() {
	// No-op for testing
}

func TestBuildDirectoryFromEvent(t *testing.T) {
	tests := []struct {
		name           string
		events         []shared.FileEvent
		expectedFiles  map[string]string // map[filename]checksum
		forceErr       bool
		errorStep      int
		expectNilDir   bool
	}{
		{
			name: "Create single file",
			events: []shared.FileEvent{
				{
					EventType:    shared.CreateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum1",
				},
			},
			expectedFiles: map[string]string{
				"file1.txt": "checksum1",
			},
		},
		{
			name: "Create and update file",
			events: []shared.FileEvent{
				{
					EventType:    shared.CreateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum1",
				},
				{
					EventType:    shared.UpdateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum2",
				},
			},
			expectedFiles: map[string]string{
				"file1.txt": "checksum2",
			},
		},
		{
			name: "Create and delete file",
			events: []shared.FileEvent{
				{
					EventType:    shared.CreateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum1",
				},
				{
					EventType: shared.DeleteEvent,
					FileName:  "file1.txt",
					FilePath:  "/path/to/file1.txt",
				},
			},
			expectedFiles: map[string]string{},
		},
		{
			name: "Multiple files with various operations",
			events: []shared.FileEvent{
				{
					EventType:    shared.CreateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum1",
				},
				{
					EventType:    shared.CreateEvent,
					FileName:     "file2.txt",
					FilePath:     "/path/to/file2.txt",
					FileChecksum: "checksum2",
				},
				{
					EventType:    shared.UpdateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum1-updated",
				},
				{
					EventType:    shared.CreateEvent,
					FileName:     "file3.txt",
					FilePath:     "/path/to/file3.txt",
					FileChecksum: "checksum3",
				},
				{
					EventType: shared.DeleteEvent,
					FileName:  "file2.txt",
					FilePath:  "/path/to/file2.txt",
				},
			},
			expectedFiles: map[string]string{
				"file1.txt": "checksum1-updated",
				"file3.txt": "checksum3",
			},
		},
		{
			name: "Update non-existent file (should add it)",
			events: []shared.FileEvent{
				{
					EventType:    shared.UpdateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum1",
				},
			},
			expectedFiles: map[string]string{
				"file1.txt": "checksum1",
			},
		},
		{
			name:      "Error during iteration",
			forceErr:  true,
			errorStep: 1,
			events: []shared.FileEvent{
				{
					EventType:    shared.CreateEvent,
					FileName:     "file1.txt",
					FilePath:     "/path/to/file1.txt",
					FileChecksum: "checksum1",
				},
				{
					EventType:    shared.CreateEvent,
					FileName:     "file2.txt",
					FilePath:     "/path/to/file2.txt",
					FileChecksum: "checksum2",
				},
			},
			expectNilDir: true,
		},
		{
			name: "Empty event collection",
			events: []shared.FileEvent{},
			expectedFiles: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockFileEventCollection{
				events:    tt.events,
				forceErr:  tt.forceErr,
				errorStep: tt.errorStep,
			}

			directory := BuildDirectoryFromEvent(mockCollection)

			if tt.expectNilDir {
				if directory != nil {
					t.Fatalf("Expected nil directory for error case, got %v", directory)
				}
				return
			}

			if directory == nil {
				t.Fatalf("Expected non-nil directory")
			}

			// Check all expected files exist with correct checksums
			for fileName, expectedChecksum := range tt.expectedFiles {
				file, exists := directory.GetFile(fileName)
				if !exists {
					t.Errorf("Expected file %s to exist in directory", fileName)
					continue
				}

				if file.GetChecksum() != expectedChecksum {
					t.Errorf("File %s has checksum %s, expected %s",
						fileName, file.GetChecksum(), expectedChecksum)
				}
			}

			// Check no unexpected files exist
			files := directory.GetFiles()
			if len(files) != len(tt.expectedFiles) {
				t.Errorf("Directory has %d files, expected %d", len(files), len(tt.expectedFiles))
			}

			for _, file := range files {
				expectedChecksum, exists := tt.expectedFiles[file.GetName()]
				if !exists {
					t.Errorf("Unexpected file in directory: %s", file.GetName())
					continue
				}

				if file.GetChecksum() != expectedChecksum {
					t.Errorf("File %s has checksum %s, expected %s",
						file.GetName(), file.GetChecksum(), expectedChecksum)
				}
			}
		})
	}
}

// TestApplyEventToDirectory tests the applyEventToDirectory function directly
func TestApplyEventToDirectory(t *testing.T) {
	tests := []struct {
		name           string
		initialFiles   map[string]*ShadowFile
		event          shared.FileEvent
		expectedFiles  map[string]string // map[filename]checksum
	}{
		{
			name: "Create new file in empty directory",
			initialFiles: map[string]*ShadowFile{},
			event: shared.FileEvent{
				EventType:    shared.CreateEvent,
				FileName:     "file1.txt",
				FilePath:     "/path/to/file1.txt",
				FileChecksum: "checksum1",
			},
			expectedFiles: map[string]string{
				"file1.txt": "checksum1",
			},
		},
		{
			name: "Update existing file",
			initialFiles: map[string]*ShadowFile{
				"file1.txt": {
					fileName: "file1.txt",
					filePath: "/path/to/file1.txt",
					checksum: "old-checksum",
				},
			},
			event: shared.FileEvent{
				EventType:    shared.UpdateEvent,
				FileName:     "file1.txt",
				FilePath:     "/path/to/file1.txt",
				FileChecksum: "new-checksum",
			},
			expectedFiles: map[string]string{
				"file1.txt": "new-checksum",
			},
		},
		{
			name: "Delete existing file",
			initialFiles: map[string]*ShadowFile{
				"file1.txt": {
					fileName: "file1.txt",
					filePath: "/path/to/file1.txt",
					checksum: "checksum1",
				},
			},
			event: shared.FileEvent{
				EventType: shared.DeleteEvent,
				FileName:  "file1.txt",
				FilePath:  "/path/to/file1.txt",
			},
			expectedFiles: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a directory with initial files
			directory := shared.Directory{
				DirectoryPath: ":memory:",
				Files:         make(map[string]shared.IFile),
			}

			// Add initial files
			for _, file := range tt.initialFiles {
				directory.AddFile(file)
			}

			// Apply the event
			applyEventToDirectory(&directory, tt.event)

			// Verify the directory state
			if len(directory.Files) != len(tt.expectedFiles) {
				t.Errorf("Directory has %d files, expected %d", len(directory.Files), len(tt.expectedFiles))
			}

			for fileName, expectedChecksum := range tt.expectedFiles {
				file, exists := directory.GetFile(fileName)
				if !exists {
					t.Errorf("Expected file %s to exist in directory", fileName)
					continue
				}

				if file.GetChecksum() != expectedChecksum {
					t.Errorf("File %s has checksum %s, expected %s",
						fileName, file.GetChecksum(), expectedChecksum)
				}
			}
		})
	}
}
