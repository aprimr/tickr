package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aprimr/tickr/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateLastLoginAndStoreRefreshToken(ctx context.Context, userID uuid.UUID, hashedToken, deviceInfo string) error

	CreateUser(ctx context.Context, req UserRegisterRequest, password_hash, otp_hash string) (uuid.UUID, error)
	CreateVenueAdmin(ctx context.Context, req VenueRegisterRequest, password_hash, otp_hash string) (uuid.UUID, error)

	GetActiveOTP(ctx context.Context, userID uuid.UUID, otpType domain.OTPType) (OTP, error)
	MarkOTPAsUsedAndVerifyUser(ctx context.Context, otpID uuid.UUID, userID uuid.UUID) error
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

// GetUserByEmail finds the user's record by email and returns the user
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, phone_number, password_hash, role, is_active, is_email_verified, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.IsEmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w", ErrUserNotFound)
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateLastLoginAndStoreRefreshToken updates the last_login_at and
// creates a new record for access token in `access_tokens` table
func (r *authRepository) UpdateLastLoginAndStoreRefreshToken(ctx context.Context, userID uuid.UUID, hashedToken, deviceInfo string) error {
	// Start the transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update last login
	updateQuery := `
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	result, err := tx.Exec(ctx, updateQuery, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w", ErrUserNotFound)
	}

	// Store refresh token
	query := `
	INSERT INTO refresh_tokens
	(user_id, hashed_token, device_info, expires_at)
	VALUES ($1, $2, $3, $4)
	`

	expiresAt := time.Now().Add(7 * 24 * time.Hour) // refresh token expires after 7 days
	_, err = tx.Exec(ctx, query, userID, hashedToken, deviceInfo, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit db transaction: %w", err)
	}

	return nil
}

// CreateUser handles user registration transaction for end users
func (r *authRepository) CreateUser(ctx context.Context, req UserRegisterRequest, password_hash, otp_hash string) (uuid.UUID, error) {
	// Begin database transction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transcation: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert a base user in users table and return its id
	var userID uuid.UUID
	userQuery := `
		INSERT INTO users (email, phone_number, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = tx.QueryRow(ctx, userQuery, req.Email, req.PhoneNumber, password_hash, RoleUser).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, ErrEmailAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Insert the details of the user in user_details table
	userDetails := `
		INSERT INTO user_details ( user_id, full_name) 
		VALUES ($1, $2)
	`
	_, err = tx.Exec(ctx, userDetails, userID, req.FullName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert user details: %w", err)
	}

	// Insert hashedOtp in the otps table
	insertOTP := `
		INSERT INTO otps (user_id, hashed_otp, type, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(ctx, insertOTP, userID, otp_hash, domain.OTPTypeAccountVerification, time.Now().Add(15*time.Minute))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert otp: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

// CreateVenueAdmin handles the registration for the venue admin
func (r *authRepository) CreateVenueAdmin(ctx context.Context, req VenueRegisterRequest, password_hash, otp_hash string) (uuid.UUID, error) {
	// Begin database transction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transcation: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert a base user in users table and return its id
	var userID uuid.UUID
	userQuery := `
		INSERT INTO users (email, phone_number, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = tx.QueryRow(ctx, userQuery, req.Email, req.PhoneNumber, password_hash, RoleVenueAdmin).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, ErrEmailAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("failed to create venue admin: %w", err)
	}

	// Insert the details of the venue in venue_details table
	userDetails := `
		INSERT INTO venue_details ( user_id, venue_name, address, city, total_screens) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, userDetails, userID, req.VenueName, req.Address, req.City, req.TotalScreens)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert venue details: %w", err)
	}

	// Insert hashedOtp in the otps table
	insertOTP := `
		INSERT INTO otps (user_id, hashed_otp, type, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(ctx, insertOTP, userID, otp_hash, domain.OTPTypeAccountVerification, time.Now().Add(15*time.Minute))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert otp: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

// GetActiveOTP fetches the newest unexpired, unused OTP record for a user
func (r *authRepository) GetActiveOTP(ctx context.Context, userID uuid.UUID, otpType domain.OTPType) (OTP, error) {
	query := `
		SELECT id, user_id, hashed_otp, type, is_used, expires_at, created_at 
		FROM otps 
		WHERE user_id = $1 AND type = $2 AND is_used = FALSE AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	var o OTP

	err := r.db.QueryRow(ctx, query, userID, otpType).Scan(
		&o.ID,
		&o.UserID,
		&o.HashedOTP,
		&o.Type,
		&o.IsUsed,
		&o.ExpiresAt,
		&o.CreatedAt,
	)
	if err != nil {
		// If no records found in the db, return invalid or expired otp error
		return OTP{}, fmt.Errorf("failed to fetch otp : %w", err)
	}

	return o, nil
}

// MarkOTPAsUsedAndVerifyUser handles the verfication of user email and mark OTP as used
func (r *authRepository) MarkOTPAsUsedAndVerifyUser(ctx context.Context, otpID uuid.UUID, userID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Mark OTP as used
	_, err = tx.Exec(ctx, `UPDATE otps SET is_used = TRUE WHERE id = $1`, otpID)
	if err != nil {
		return fmt.Errorf("failed to update otp status: %w", err)
	}

	// Mark user email as verified
	_, err = tx.Exec(ctx, `UPDATE users SET is_email_verified = TRUE WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to verify user email status: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
