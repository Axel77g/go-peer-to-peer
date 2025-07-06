package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type IEncrypter interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type AESEncrypter struct {
	Key []byte // AES key (16, 24, or 32 bytes for AES-128, AES-192, or AES-256)
}

func NewAESEncrypter(key []byte) (*AESEncrypter, error) {
	// Validate key length (must be 16, 24, or 32 bytes)
	switch len(key) {
	case 16, 24, 32:
		return &AESEncrypter{Key: key}, nil
	default:
		return nil, errors.New("invalid key size: must be 16, 24, or 32 bytes")
	}
}

func (e *AESEncrypter) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, 12) // 12 bytes for GCM
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, iv, data, nil)
	return append(iv, ciphertext...), nil
}

func (e *AESEncrypter) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ivSize := 12 // 12 bytes for GCM
	if len(data) < ivSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:ivSize]
	ciphertext := data[ivSize:]

	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// Example usage:

/* // Create a new encrypter with a 32-byte key (AES-256)
key := make([]byte, 32)
// In a real app, use a secure key generation method
rand.Read(key)

encrypter, err := encryption.NewAESEncrypter(key)
if err != nil {
    // Handle error
}

// Encrypt data
encrypted, err := encrypter.Encrypt([]byte("sensitive data"))
if err != nil {
    // Handle error
}

// Decrypt data
decrypted, err := encrypter.Decrypt(encrypted)
if err != nil {
    // Handle error
} */