package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// CheckPasswordHash compares a plaintext password with a bcrypt hash.
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashSHA256 creates a SHA-256 hash of the input data.
func HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

type EncryptorManager struct {
	gcm cipher.AEAD
}

func NewEncryptorManager(key []byte) (*EncryptorManager, error) {

	if len(key) != 32 {
		return nil, errors.New("Encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	return &EncryptorManager{gcm: gcm}, nil
}

func (em *EncryptorManager) EncryptSecret(secret string) (string, error) {
	nonce := make([]byte, em.gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := em.gcm.Seal(nonce, nonce, []byte(secret), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (em *EncryptorManager) DecryptSecret(encryptedText string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)

	if err != nil {
		return "", err
	}

	nonce, cipherText := cipherText[:em.gcm.NonceSize()], cipherText[em.gcm.NonceSize():]

	plainText, err := em.gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}
