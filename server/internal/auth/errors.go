package auth

import "errors"

var (
	ErrEmailNotVerified   = errors.New("email is not verified, please verify it first")
	ErrAccountDeactivated = errors.New("account deactivated, please contact support")

	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidOrExpiredOTP = errors.New("invalid or expired verification code")

	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrFailedToCreateToken = errors.New("failed to create jwt token")
)
