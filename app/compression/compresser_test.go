package compression

import (
	"bytes"
	"compress/gzip"
	"reflect"
	"testing"
)

func TestNewGzipCompresser(t *testing.T) {
	tests := []struct {
		name    string
		level   int
		wantErr bool
	}{
		{"NoCompression", gzip.NoCompression, false},
		{"BestSpeed", gzip.BestSpeed, false},
		{"BestCompression", gzip.BestCompression, false},
		{"DefaultCompression", gzip.DefaultCompression, false},
		{"HuffmanOnly", gzip.HuffmanOnly, false},
		{"CustomValidLevel", 5, false},
		{"InvalidLevel", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGzipCompresser(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGzipCompresser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Level != tt.level {
				t.Errorf("NewGzipCompresser() got level = %v, want %v", got.Level, tt.level)
			}
		})
	}
}

func TestGzipCompresser_Compress(t *testing.T) {
	testData := []byte("This is some test data for compression")
	
	compresser, err := NewGzipCompresser(gzip.DefaultCompression)
	if err != nil {
		t.Fatalf("Failed to create compresser: %v", err)
	}

	compressed, err := compresser.Compress(testData)
	if err != nil {
		t.Fatalf("Compress() error = %v", err)
	}

	// Verify compressed data is not empty and different from original
	if len(compressed) == 0 {
		t.Error("Compress() returned empty data")
	}

	if bytes.Equal(compressed, testData) {
		t.Error("Compress() did not change data, compression failed")
	}

	// Verify compressed data is smaller (should be true for text data)
	if len(compressed) >= len(testData) {
		t.Logf("Note: Compressed data is not smaller than original (%d >= %d bytes)", 
			len(compressed), len(testData))
	}
}

func TestGzipCompresser_Decompress(t *testing.T) {
	testData := []byte("This is some test data for compression and decompression")
	
	compresser, err := NewGzipCompresser(gzip.DefaultCompression)
	if err != nil {
		t.Fatalf("Failed to create compresser: %v", err)
	}

	compressed, err := compresser.Compress(testData)
	if err != nil {
		t.Fatalf("Compress() error = %v", err)
	}

	decompressed, err := compresser.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress() error = %v", err)
	}

	// Verify decompressed data matches original
	if !bytes.Equal(decompressed, testData) {
		t.Errorf("Decompress() = %v, want %v", decompressed, testData)
	}
}

func TestGzipCompresser_CompressEmpty(t *testing.T) {
	emptyData := []byte{}
	
	compresser, err := NewGzipCompresser(gzip.DefaultCompression)
	if err != nil {
		t.Fatalf("Failed to create compresser: %v", err)
	}

	compressed, err := compresser.Compress(emptyData)
	if err != nil {
		t.Fatalf("Compress() error = %v", err)
	}

	// Empty data should still produce some compressed output (gzip headers)
	if len(compressed) == 0 {
		t.Error("Compress() of empty data returned empty result")
	}

	decompressed, err := compresser.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress() error = %v", err)
	}

	// Verify decompressed data matches original empty data
	if !bytes.Equal(decompressed, emptyData) {
		t.Errorf("Decompress() = %v, want empty slice", decompressed)
	}
}

func TestGzipCompresser_DecompressInvalid(t *testing.T) {
	invalidData := []byte("This is not valid compressed data")
	
	compresser, _ := NewGzipCompresser(gzip.DefaultCompression)
	_, err := compresser.Decompress(invalidData)
	
	if err == nil {
		t.Error("Decompress() did not return error for invalid data")
	}
}

func TestGzipCompresser_DecompressEmpty(t *testing.T) {
	emptyData := []byte{}
	
	compresser, _ := NewGzipCompresser(gzip.DefaultCompression)
	_, err := compresser.Decompress(emptyData)
	
	if err == nil {
		t.Error("Decompress() did not return error for empty data")
	}
}

func TestGzipCompresser_RoundTrip(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{"Small text", []byte("Hello, world!")},
		{"Empty", []byte{}},
		{"Binary data", bytes.Repeat([]byte{0xFF, 0x00, 0xAA, 0x55}, 100)},
		{"Large text", bytes.Repeat([]byte("Lorem ipsum dolor sit amet "), 1000)},
	}

	for _, level := range []int{
		gzip.NoCompression,
		gzip.BestSpeed,
		gzip.DefaultCompression,
		gzip.BestCompression,
	} {
		compresser, err := NewGzipCompresser(level)
		if err != nil {
			t.Fatalf("Failed to create compresser with level %d: %v", level, err)
		}

		for _, tc := range testCases {
			t.Run(tc.name+"_Level_"+string(rune(level+'0')), func(t *testing.T) {
				compressed, err := compresser.Compress(tc.data)
				if err != nil {
					t.Fatalf("Compress() error = %v", err)
				}

				decompressed, err := compresser.Decompress(compressed)
				if err != nil {
					t.Fatalf("Decompress() error = %v", err)
				}

				if !reflect.DeepEqual(decompressed, tc.data) {
					t.Errorf("Round trip failed: got %v, want %v", decompressed, tc.data)
				}
			})
		}
	}
}
