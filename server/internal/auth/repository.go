package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, req UserRegisterRequest, password_hash string) (uuid.UUID, error)
	CreateVenueAdmin(ctx context.Context, req VenueRegisterRequest, password_hash string) (uuid.UUID, error)
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

// CreateUser handles user registration transaction for end users
func (r *authRepository) CreateUser(ctx context.Context, req UserRegisterRequest, password_hash string) (uuid.UUID, error) {
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
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Insert the details of the user in user_details table
	userDetails := `
		INSERT INTO user_details ( user_id, full_name) 
		VALUES ($1, $2)
	`
	_, err = r.db.Exec(ctx, userDetails, userID, req.FullName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert user details: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

// CreateVenueAdmin handles the registration for the venue admin
func (r *authRepository) CreateVenueAdmin(ctx context.Context, req VenueRegisterRequest, password_hash string) (uuid.UUID, error) {
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
		return uuid.Nil, fmt.Errorf("failed to create venue admin: %w", err)
	}

	// Insert the details of the venue in venue_details table
	userDetails := `
		INSERT INTO venue_details ( user_id, venue_name, address, city, total_screens) 
		VALUES ($1, $2)
	`
	_, err = r.db.Exec(ctx, userDetails, userID, req.VenueName, req.Address, req.City, req.TotalScreens)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert venue details: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}
