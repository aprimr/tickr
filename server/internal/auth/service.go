package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/aprimr/tickr/internal/domain"
	"github.com/aprimr/tickr/internal/email"
	"github.com/aprimr/tickr/internal/utils/hash"
	"github.com/aprimr/tickr/internal/utils/jwt"
	"github.com/aprimr/tickr/internal/utils/otp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthService interface {
	Login(ctx context.Context, req LoginRequest, deviceInfo string) (string, string, error)

	RegisterUser(ctx context.Context, req UserRegisterRequest) (uuid.UUID, error)
	RegisterVenueAdmin(ctx context.Context, req VenueRegisterRequest) (uuid.UUID, error)

	VerifyUserAccount(ctx context.Context, req VerifyAccountRequest) error
}

type authService struct {
	repo   AuthRepository
	mailer email.EmailService
}

func NewAuthService(repo AuthRepository, mailer email.EmailService) AuthService {
	return &authService{
		repo:   repo,
		mailer: mailer,
	}
}

// Login handles authentication for users, venue admins and superadmins
// and returns the access token, refresh token and error
func (s *authService) Login(ctx context.Context, req LoginRequest, deviceInfo string) (string, string, error) {
	// Call repository to find user in db
	userDetail, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", "", ErrInvalidCredentials
		}

		return "", "", fmt.Errorf("auth service error: %w", err)
	}

	// Check if user is active/not banned
	if !userDetail.IsActive {
		return "", "", ErrAccountDeactivated
	}

	// Check if email is verified
	if !userDetail.IsEmailVerified {
		return "", "", ErrEmailNotVerified
	}

	// Compare user's password against stored password on db
	err = CheckPassword(userDetail.PasswordHash, req.Password)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate access token
	accessToken, err := jwt.GenerateAccessToken(userDetail.ID, string(userDetail.Role))
	if err != nil {
		return "", "", ErrFailedToCreateToken
	}

	// Generate refresh token
	refreshToken, err := jwt.GenerateRefreshToken(userDetail.ID, string(userDetail.Role))
	if err != nil {
		return "", "", ErrFailedToCreateToken
	}

	// Hash refresh token
	hashedRefreshToken, err := hash.String(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash refresh token: %w", err)
	}

	// Update last login and store refresh token
	err = s.repo.UpdateLastLoginAndStoreRefreshToken(ctx, userDetail.ID, hashedRefreshToken, deviceInfo)
	if err != nil {
		return "", "", fmt.Errorf("failed to update last login and store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// RegisterCustomer handles user signup
func (s *authService) RegisterUser(ctx context.Context, req UserRegisterRequest) (uuid.UUID, error) {
	// Generate OTP
	otpString, err := otp.GenerateOTP()
	if err != nil {
		return uuid.Nil, err
	}

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return uuid.Nil, err
	}

	// Hash OTP
	hashedOTP, err := hash.String(otpString)
	if err != nil {
		return uuid.Nil, err
	}

	// Call repository to create user
	userID, err := s.repo.CreateUser(ctx, req, hashedPassword, hashedOTP)
	if err != nil {
		return uuid.Nil, err
	}

	// Send mail, If user is created successfully
	_ = s.mailer.SendAccountVerificationEmail(req.Email, req.FullName, otpString)

	return userID, nil
}

// RegisterVenueAdmin handles venue signup
func (s *authService) RegisterVenueAdmin(ctx context.Context, req VenueRegisterRequest) (uuid.UUID, error) {
	// Generate OTP
	otpString, err := otp.GenerateOTP()
	if err != nil {
		return uuid.Nil, err
	}

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return uuid.Nil, err
	}

	// Hash OTP
	hashedOTP, err := hash.String(otpString)
	if err != nil {
		return uuid.Nil, err
	}

	// Call repository
	userID, err := s.repo.CreateVenueAdmin(ctx, req, hashedPassword, hashedOTP)
	if err != nil {
		return uuid.Nil, err
	}

	// Send mail, If user is created successfully
	_ = s.mailer.SendAccountVerificationEmail(req.Email, req.VenueName, otpString)

	return userID, nil
}

func (s *authService) VerifyUserAccount(ctx context.Context, req VerifyAccountRequest) error {
	// Fetch the active OTP from database
	otpRecord, err := s.repo.GetActiveOTP(ctx, req.UserID, domain.OTPTypeAccountVerification)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidOrExpiredOTP
		}
		return err
	}

	// Compare plain OTP with the hashedOTP
	match := hash.CheckString(req.OTP, otpRecord.HashedOTP)
	if !match {
		return ErrInvalidOrExpiredOTP
	}

	// Mark OTP as used and verify user account
	err = s.repo.MarkOTPAsUsedAndVerifyUser(ctx, otpRecord.ID, req.UserID)
	if err != nil {
		return err
	}

	return nil
}
