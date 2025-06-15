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
	// Create a new cipher block from the key
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore, it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	// Fill iv with random bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Use CFB mode for encryption
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func (e *AESEncrypter) Decrypt(data []byte) ([]byte, error) {
	// Check if data is long enough to contain IV and at least some ciphertext
	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	// Extract the IV from the beginning of the ciphertext
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	// Create the plaintext with the same length as ciphertext
	plaintext := make([]byte, len(ciphertext))

	// Use CFB mode for decryption
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

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