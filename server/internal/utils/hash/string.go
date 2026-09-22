package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString takes a plain text string and returns its bcrypt hash
func String(input string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckHash compares a plain text string with a hashed string
func CheckString(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}
