package auth

import (
	"context"

	"github.com/aprimr/tickr/internal/email"
	"github.com/aprimr/tickr/internal/utils/hash"
	"github.com/aprimr/tickr/internal/utils/otp"
	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(ctx context.Context, req UserRegisterRequest) (uuid.UUID, error)
	RegisterVenueAdmin(ctx context.Context, req VenueRegisterRequest) (uuid.UUID, error)
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
