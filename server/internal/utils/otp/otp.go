package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP creates a 6-digit numeric OTP string
func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
