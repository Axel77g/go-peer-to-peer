package filetransfering

import (
	"strconv"
)

type Chunk struct {
	Pos    uint16
	Size   uint32
	Buffer []byte
}

func (chunk *Chunk) GetChunckPayload() []byte {
	prefix := []byte("CHUNK_" + strconv.Itoa(int(chunk.Pos)) + ",")
	result := make([]byte, int(chunk.Size))
	result = append(result, prefix...)
	result = append(result, chunk.Buffer...)
	return result
}

type File struct {
	ID        uint16
	Size      uint64
	Name      string
	ChunkSize uint32
	Chunks    []Chunk
}

func NewFile(id uint16,
	size uint64,
	name string,
	chunckSize uint32,
	buffer []byte) *File {
	chunks := make([]Chunk, 0)
	for i := uint64(0); i < size/uint64(chunckSize); i++ {
		chunks = append(chunks, Chunk{
			Pos:    uint16(i),
			Size:   chunckSize,
			Buffer: buffer[i*uint64(chunckSize) : (i+1)*uint64(chunckSize)],
		})
	}

	if size%uint64(chunckSize) != 0 {
		chunks = append(chunks, Chunk{
			Pos:    uint16(size / uint64(chunckSize)),
			Size:   uint32(size % uint64(chunckSize)),
			Buffer: buffer[size/uint64(chunckSize)*uint64(chunckSize):],
		})
	}

	return &File{
		ID:        id,
		Size:      size,
		Name:      name,
		ChunkSize: chunckSize,
		Chunks:    chunks,
	}
}

func (file *File) GetChunk(pos int) Chunk {
	return file.Chunks[pos]
}

func ReadFile(filePath string) *File {
	//Read file

	return NewFile(1, 7, "file", 1, []byte("bonjour"))
}
