package encryption

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestAESEncrypter_Creation(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{
			name:    "Valid AES-128 Key",
			keySize: 16,
			wantErr: false,
		},
		{
			name:    "Valid AES-192 Key",
			keySize: 24,
			wantErr: false,
		},
		{
			name:    "Valid AES-256 Key",
			keySize: 32,
			wantErr: false,
		},
		{
			name:    "Invalid Key Size",
			keySize: 18,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)
			_, err := rand.Read(key)
			if err != nil {
				t.Fatalf("Failed to generate random key: %v", err)
			}

			encrypter, err := NewAESEncrypter(key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESEncrypter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && encrypter == nil {
				t.Error("Expected valid encrypter, got nil")
			}
		})
	}
}

func TestAESEncrypter_EncryptDecrypt(t *testing.T) {
	// Create different test data sizes
	testData := []struct {
		name string
		data []byte
	}{
		{name: "Empty data", data: []byte("")},
		{name: "Short text", data: []byte("Hello World")},
		{name: "Longer text", data: []byte("This is a longer text that will be encrypted and then decrypted")},
		{name: "Binary data", data: func() []byte {
			data := make([]byte, 1024)
			rand.Read(data)
			return data
		}()},
	}

	// Test with different key sizes
	keySizes := []int{16, 24, 32} // AES-128, AES-192, AES-256

	for _, keySize := range keySizes {
		key := make([]byte, keySize)
		_, err := rand.Read(key)
		if err != nil {
			t.Fatalf("Failed to generate random key: %v", err)
		}

		encrypter, err := NewAESEncrypter(key)
		if err != nil {
			t.Fatalf("Failed to create encrypter with key size %d: %v", keySize, err)
		}

		for _, tt := range testData {
			bitSize := keySize * 8
			t.Run(fmt.Sprintf("%s with AES-%d", tt.name, bitSize), func(t *testing.T) {
				// Encrypt the data
				encrypted, err := encrypter.Encrypt(tt.data)
				if err != nil {
					t.Fatalf("Encryption failed: %v", err)
				}

				// Make sure encrypted data is different from original
				if len(tt.data) > 0 && bytes.Equal(encrypted, tt.data) {
					t.Error("Encrypted data is identical to original data")
				}

				// Decrypt the data
				decrypted, err := encrypter.Decrypt(encrypted)
				if err != nil {
					t.Fatalf("Decryption failed: %v", err)
				}

				// Verify that the decrypted data matches the original
				if !bytes.Equal(decrypted, tt.data) {
					t.Errorf("Decrypted data does not match original. Original: %v, Decrypted: %v", tt.data, decrypted)
				}
			})
		}
	}
}

func TestAESEncrypter_DecryptionFailure(t *testing.T) {
	// Generate a valid key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	encrypter, err := NewAESEncrypter(key)
	if err != nil {
		t.Fatalf("Failed to create encrypter: %v", err)
	}

	// Test cases for decryption failures
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "Too short data",
			data: []byte("too short"),
		},
		{
			name: "Tampered data",
			data: func() []byte {
				originalData := []byte("Original data")
				encrypted, _ := encrypter.Encrypt(originalData)
				// Tamper with the ciphertext portion
				if len(encrypted) > 15 {
					encrypted[15] ^= 0xFF // Flip bits in the ciphertext
				}
				return encrypted
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encrypter.Decrypt(tt.data)
			if err == nil {
				t.Error("Expected error during decryption but got nil")
			}
		})
	}
}
