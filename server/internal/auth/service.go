package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aprimr/tickr/internal/cache"
	"github.com/aprimr/tickr/internal/domain"
	"github.com/aprimr/tickr/internal/email"
	"github.com/aprimr/tickr/internal/utils/hash"
	"github.com/aprimr/tickr/internal/utils/jwt"
	"github.com/aprimr/tickr/internal/utils/otp"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req LoginRequest, deviceInfo string) (string, string, error)

	RegisterUser(ctx context.Context, req UserRegisterRequest) (uuid.UUID, error)
	RegisterVenueAdmin(ctx context.Context, req VenueRegisterRequest) (uuid.UUID, error)

	VerifyUserAccount(ctx context.Context, req VerifyAccountRequest) error
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error

	RotateToken(ctx context.Context, req RotateTokenRequest, deviceInfo string) (string, string, error)

	Logout(ctx context.Context, refreshToken string) error
	LogoutAllDevices(ctx context.Context, userID uuid.UUID) error
	GetActiveSessions(ctx context.Context, userID uuid.UUID, refreshToken string) ([]SessionResponse, error)
	DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
}

type authService struct {
	repo   AuthRepository
	mailer email.EmailService
	cache  *cache.CacheClient
}

func NewAuthService(repo AuthRepository, mailer email.EmailService, cache *cache.CacheClient) AuthService {
	return &authService{
		repo:   repo,
		mailer: mailer,
		cache:  cache,
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

	// Call repository to create user
	userID, err := s.repo.CreateUser(ctx, req, hashedPassword)
	if err != nil {
		return uuid.Nil, err
	}

	// Store OTP in redis
	err = s.cache.StoreOTP(ctx, domain.OTPTypeAccountVerification, userID.String(), otpString)
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

	// Call repository
	userID, err := s.repo.CreateVenueAdmin(ctx, req, hashedPassword)
	if err != nil {
		return uuid.Nil, err
	}

	// Store OTP in redis
	err = s.cache.StoreOTP(ctx, domain.OTPTypeAccountVerification, userID.String(), otpString)
	if err != nil {
		return uuid.Nil, err
	}

	// Send mail, If user is created successfully
	_ = s.mailer.SendAccountVerificationEmail(req.Email, req.VenueName, otpString)

	return userID, nil
}

// VerifyUserAccount handles the user's email verification
func (s *authService) VerifyUserAccount(ctx context.Context, req VerifyAccountRequest) error {
	// Get OTP from redis
	err := s.cache.VerifyOTP(ctx, domain.OTPTypeAccountVerification, req.UserID.String(), req.OTP)
	if err != nil {
		return err
	}

	// Mark OTP as used and verify user account
	err = s.repo.VerifyUserAccount(ctx, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

// ForgotPassword handles the user's account password reset request
func (s *authService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// Fetch user from db
	userDetail, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Return nil if the err is user not found
		if errors.Is(err, ErrUserNotFound) {
			return nil
		}

		// Else return actual error
		return err
	}

	// If any active otp, return error, ignore the returned error
	exists, err := s.cache.OTPExists(ctx, domain.OTPTypeForgotPassword, req.Email)
	if err != nil {
		return fmt.Errorf("failed to check active otp: %w", err)
	}
	if exists {
		return fmt.Errorf("%w", ErrActiveOTPAlreadyExists)
	}

	// Generate OTP
	otp, err := otp.GenerateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate otp: %w", err)
	}

	// Store OTP in redis
	err = s.cache.StoreOTP(ctx, domain.OTPTypeForgotPassword, userDetail.Email, otp)
	if err != nil {
		return fmt.Errorf("failed to store otp: %w", err)
	}

	// Send email
	err = s.mailer.SendForgotPasswordEmail(req.Email, "there", otp)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// ResetPassword updates the user's password, mark user verified and also revokes all the sessions
func (s *authService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// Fetch user from db
	userDetail, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// If user is not found, return `invalid or expired otp`
		if errors.Is(err, ErrUserNotFound) {
			return fmt.Errorf("%w", ErrInvalidOrExpiredOTP)
		}

		// Else return actual error
		return err
	}

	// Verify OTP from cache
	if err := s.cache.VerifyOTP(ctx, domain.OTPTypeForgotPassword, userDetail.Email, req.OTP); err != nil {
		if errors.Is(err, cache.ErrInvalidOTP) {
			return fmt.Errorf("%w", ErrInvalidOrExpiredOTP)
		}
		return err
	}

	// Hash new password
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Set new password
	err = s.repo.ResetPasswordAndRevokeSessions(ctx, userDetail.ID, hashedPassword)
	if err != nil {
		return err
	}

	// Send password changed email
	err = s.mailer.SendPasswordResetSuccessEmail(userDetail.Email, "there")
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// RotateToken verifies the refresh token, replaces it in the database,
// returns access token, refresh token and error
func (s *authService) RotateToken(ctx context.Context, req RotateTokenRequest, deviceInfo string) (string, string, error) {
	// Verify jwt signature
	claims, err := jwt.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return "", "", jwt.ErrExpiredToken
		}

		return "", "", jwt.ErrInvalidToken
	}

	// Hash incomming refresh token
	hashedClientToken, err := hash.String(req.RefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash refresh token: %w", err)
	}

	// Get refresh token from db
	refreshToken, err := s.repo.GetRefreshTokenByItsHash(ctx, hashedClientToken)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) {
			return "", "", domain.ErrRefreshTokenNotFound
		}

		return "", "", fmt.Errorf("failed to get refresh token: %w", err)

	}

	// Verify if both tokens belong to same user
	if claims.UserID != refreshToken.UserID {
		return "", "", jwt.ErrInvalidToken
	}

	// Verify token expiration
	if time.Now().After(refreshToken.ExpiresAt) {
		return "", "", jwt.ErrExpiredToken
	}

	// Generate new access and refresh token
	newAccessToken, err := jwt.GenerateAccessToken(claims.UserID, claims.UserRole)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefreshToken, err := jwt.GenerateRefreshToken(claims.UserID, claims.UserRole)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash new refresh token
	newRefreshHashed, err := hash.String(newRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash refresh token: %w", err)
	}

	// Replace refresh token
	err = s.repo.DeleteOldAndCreateNewRefreshToken(ctx, refreshToken.ID, claims.UserID, newRefreshHashed, deviceInfo)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout deletes the active refresh token session
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash, err := hash.String(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to hash refresh token: %w", err)
	}

	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return err
	}

	return nil
}

// LogoutAllDevices revokes all active sessions for a user.
func (s *authService) LogoutAllDevices(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.DeleteAllUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to logout from all devices: %w", err)
	}
	return nil
}

// GetActiveSessions fetches the list of active session for a user and flag the current session
func (s *authService) GetActiveSessions(ctx context.Context, userID uuid.UUID, refreshToken string) ([]SessionResponse, error) {
	sessions, err := s.repo.GetActiveSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sessionResponse []SessionResponse
	for _, session := range *sessions {

		isCurrentSession := hash.CheckString(refreshToken, session.HashedToken)

		resp := SessionResponse{
			SessionID:  session.ID.String(),
			DeviceInfo: session.DeviceInfo,
			IsCurrent:  isCurrentSession,
			CreatedAt:  session.CreatedAt,
		}

		sessionResponse = append(sessionResponse, resp)
	}

	return sessionResponse, nil
}

// DeleteSession revokes a specific session for a user
func (s *authService) DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	err := s.repo.DeleteSessionByID(ctx, userID, sessionID)
	if err != nil {
		return err
	}

	return nil
}
