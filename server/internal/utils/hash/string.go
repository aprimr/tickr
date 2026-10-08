package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Load on startup
var hmacSecret = []byte(os.Getenv("HMAC_HASH_SECRET"))

// String takes a plain text string and returns its SHA-256 hash
// It will always return the exact same output, for a given string
// (for jwt tokens)
func String(input string) (string, error) {
	hashBytes := sha256.Sum256([]byte(input))

	return hex.EncodeToString(hashBytes[:]), nil
}

// CheckString compares a plain text string with a hashed string
func CheckString(input, hash string) bool {
	incomingHash, err := String(input)
	if err != nil {
		return false
	}

	return incomingHash == hash
}

// BcryptString takes a short plain string and returns its bcrypt hash.
// (for sensitive data)
func BcryptString(input string) (string, error) {
	if len(input) > 72 {
		return "", ErrStringTooLong
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckBcryptString compares a plaintext string with a bcrypt hash.
func CheckBcryptString(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))

	return err == nil
}

// HMACString takes a short plain string and returns HMAC-SHA256 hex hash.
// (for fast hashing)
func HMACString(input string) string {
	h := hmac.New(sha256.New, []byte(hmacSecret))
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// CheckHMACString compares a plaintext string with a HMAC hash.
func CheckHMACString(input string, hash string) bool {
	actualHash := HMACString(input)
	return hmac.Equal([]byte(actualHash), []byte(hash))
}
