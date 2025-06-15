package compression

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
)

type ICompresser interface {
    Compress(data []byte) ([]byte, error)
    Decompress(data []byte) ([]byte, error)
}

type GzipCompresser struct {
    Level int // Compression level (gzip.NoCompression, gzip.BestSpeed, gzip.BestCompression, gzip.DefaultCompression)
}

func NewGzipCompresser(level int) (*GzipCompresser, error) {
    // Validate compression level
    switch level {
    case gzip.NoCompression, gzip.BestSpeed, gzip.BestCompression, gzip.DefaultCompression:
        return &GzipCompresser{Level: level}, nil
    case gzip.HuffmanOnly:
        return &GzipCompresser{Level: level}, nil
    default:
        if level >= gzip.BestSpeed && level <= gzip.BestCompression {
            return &GzipCompresser{Level: level}, nil
        }
        return nil, errors.New("invalid compression level")
    }
}

func (c *GzipCompresser) Compress(data []byte) ([]byte, error) {
    var buf bytes.Buffer
    
    // Create a new gzip writer with the specified compression level
    writer, err := gzip.NewWriterLevel(&buf, c.Level)
    if err != nil {
        return nil, err
    }
    
    // Write data to the gzip writer
    _, err = writer.Write(data)
    if err != nil {
        writer.Close()
        return nil, err
    }
    
    // Close the writer to flush any remaining data
    err = writer.Close()
    if err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

func (c *GzipCompresser) Decompress(data []byte) ([]byte, error) {
    // Check if data is not empty
    if len(data) == 0 {
        return nil, errors.New("compressed data is empty")
    }
    
    // Create a reader from the compressed data
    reader := bytes.NewReader(data)
    
    // Create a gzip reader
    gzipReader, err := gzip.NewReader(reader)
    if err != nil {
        return nil, err
    }
    defer gzipReader.Close()
    
    // Read all decompressed data
    decompressed, err := io.ReadAll(gzipReader)
    if err != nil {
        return nil, err
    }
    
    return decompressed, nil
}

// Example usage:
/* // Create a new compresser with default compression
compresser, err := compression.NewGzipCompresser(gzip.DefaultCompression)
if err != nil {
    // Handle error
}

// Compress data
compressed, err := compresser.Compress([]byte("This is some data to compress"))
if err != nil {
    // Handle error
}

// Decompress data
decompressed, err := compresser.Decompress(compressed)
if err != nil {
    // Handle error
} */