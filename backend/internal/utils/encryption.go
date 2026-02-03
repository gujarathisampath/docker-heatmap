package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"docker-heatmap/internal/config"
)

var (
	ErrInvalidKey    = errors.New("encryption key must be exactly 32 bytes for AES-256")
	ErrDecryptFailed = errors.New("failed to decrypt data (wrong key or corrupted data)")
	ErrInvalidIV     = errors.New("invalid initialization vector")
)

// Encrypt encrypts plaintext using AES-256-GCM (Industry Standard)
// Returns the base64-encoded ciphertext and IV (nonce)
func Encrypt(plaintext string) (ciphertext, iv string, err error) {
	key := []byte(config.AppConfig.EncryptionKey)
	if len(key) != 32 {
		return "", "", ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	// GCM provides both confidentiality and authenticity (AEAD)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	// Never reuse a nonce with the same key
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	// Seal handles the encryption and appends an authentication tag
	encrypted := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(encrypted),
		base64.StdEncoding.EncodeToString(nonce),
		nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM
func Decrypt(ciphertext, iv string) (string, error) {
	key := []byte(config.AppConfig.EncryptionKey)
	if len(key) != 32 {
		return "", ErrInvalidKey
	}

	encryptedData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	nonce, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", ErrInvalidIV
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Open decrypts and authenticates the data
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}

	return string(plaintext), nil
}

// GenerateRandomString generates a cryptographically secure random string
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
