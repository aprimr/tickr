package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	ErrExpiredToken = errors.New("jwt token is expired")
	ErrInvalidToken = errors.New("invalid jwt token")
)

type RefreshToken struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	HashedToken string    `db:"hashed_token"`
	DeviceInfo  string    `db:"device_info"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
}
