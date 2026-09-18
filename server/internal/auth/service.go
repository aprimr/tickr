package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(ctx context.Context, req UserRegisterRequest) (uuid.UUID, error)
	RegisterVenueAdmin(ctx context.Context, req VenueRegisterRequest) (uuid.UUID, error)
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{repo: repo}
}

// RegisterCustomer handles user signup
func (s *authService) RegisterUser(ctx context.Context, req UserRegisterRequest) (uuid.UUID, error) {
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

	return userID, nil
}

// RegisterVenueAdmin handles venue signup
func (s *authService) RegisterVenueAdmin(ctx context.Context, req VenueRegisterRequest) (uuid.UUID, error) {
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

	return userID, nil
}
