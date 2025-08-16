package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

// GenerateUUIDToken
func GenerateUUIDToken() string {
	return uuid.New().String()
}

// Generate random token
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
