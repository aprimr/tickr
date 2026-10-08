package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aprimr/tickr/internal/domain"
	"github.com/aprimr/tickr/internal/utils/hash"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidOTP = errors.New("invalid or expired otp")
)

// genOTPKey is helper to build a consistent key format
func genOTPKey(otpType domain.OTPType, identifier string) string {
	return fmt.Sprintf("otp:%s:%s", otpType, identifier)
}

// StoreOTP saves the hashed otp with TTL of 15 minutes
func (c *CacheClient) StoreOTP(ctx context.Context, otpType domain.OTPType, identifier, otp string) error {
	// Hash otp
	hashedOTP := hash.HMACString(otp)

	// Generate key
	key := genOTPKey(otpType, identifier)

	result := c.Client.Set(ctx, key, hashedOTP, 15*time.Minute)
	return result.Err()
}

// VerifyOTP checks the stored OTP hash in Redis and verifies the OTP code
func (c *CacheClient) VerifyOTP(ctx context.Context, otpType domain.OTPType, identifier, otp string) error {
	// generate key
	key := genOTPKey(otpType, identifier)

	// get stored otp
	storedHash, err := c.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return ErrInvalidOTP
	} else if err != nil {
		return fmt.Errorf("failed to get hashedOTP: %w", err)
	}

	if !hash.CheckHMACString(otp, storedHash) {
		return ErrInvalidOTP
	}

	// delete key if matched
	err = c.Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete otp: %w", err)
	}

	return nil
}

// OTPExists checks whether an unexpired OTP exists for the identifier without deleting it
func (c *CacheClient) OTPExists(ctx context.Context, otpType domain.OTPType, identifier string) (bool, error) {
	key := genOTPKey(otpType, identifier)
	n, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check otp existence: %w", err)
	}
	return n > 0, nil
}
