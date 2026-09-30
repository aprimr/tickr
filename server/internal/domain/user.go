package domain

type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleVenueAdmin UserRole = "venue_admin"
	RoleSuperAdmin UserRole = "super_admin"
)
