package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string
type VenueStatus string

const (
	RoleUser       UserRole = "user"
	RoleVenueAdmin UserRole = "venue_admin"
	RoleSuperAdmin UserRole = "super_admin"
)

const (
	StatusPending  VenueStatus = "pending"
	StatusApproved VenueStatus = "approved"
	StatusRejected VenueStatus = "rejected"
)

type User struct {
	ID              uuid.UUID  `db:"id"`
	Email           string     `db:"email"`
	PhoneNumber     string     `db:"phone_number"`
	PasswordHash    string     `db:"password_hash"`
	Role            UserRole   `db:"role"`
	IsActive        bool       `db:"is_active"`
	IsEmailVerified bool       `db:"is_email_verified"`
	IsPhoneVerified bool       `db:"is_phone_verified"`
	LastLoginAt     *time.Time `db:"last_login_at"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

type UserDetails struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	FullName  string    `db:"full_name"`
	UpdatedAt time.Time `db:"updated_at"`
}

type VenueDetails struct {
	ID             uuid.UUID   `db:"id"`
	UserID         uuid.UUID   `db:"user_id"`
	VenueName      string      `db:"venue_name"`
	Address        string      `db:"address"`
	City           string      `db:"city"`
	TotalScreens   int         `db:"total_screens"`
	VenueStatus    VenueStatus `db:"venue_status"`
	ApprovedBy     *uuid.UUID  `db:"approved_by"`
	ApprovedAt     *time.Time  `db:"approved_at"`
	RejectedBy     *uuid.UUID  `db:"rejected_by"`
	RejectedAt     *time.Time  `db:"rejected_at"`
	VenueCreatedAt time.Time   `db:"venue_created_at"`
	VenueUpdatedAt time.Time   `db:"venue_updated_at"`
}

type AdminDetails struct {
	ID             uuid.UUID `db:"id"`
	UserID         uuid.UUID `db:"user_id"`
	FullName       string    `db:"full_name"`
	IsRootUser     bool      `db:"is_root_user"`
	AdminCreatedAt time.Time `db:"admin_created_at"`
	AdminUpdatedAt time.Time `db:"admin_updated_at"`
}
