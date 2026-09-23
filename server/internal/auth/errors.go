package auth

import "errors"

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidOrExpiredOTP = errors.New("invalid or expired verification code")
	ErrUserNotFound        = errors.New("user not found")
)
