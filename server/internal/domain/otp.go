package domain

import (
	"time"

	"github.com/google/uuid"
)

type OTPType string

const (
	OTPTypeAccountVerification OTPType = "account_verification"
	OTPTypeForgotPassword      OTPType = "forgot_password"
)

type OTP struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	HashedOTP string    `db:"hashed_otp"`
	Type      OTPType   `db:"type"`
	IsUsed    bool      `db:"is_used"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
