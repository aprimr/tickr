package domain

type VenueStatus string

const (
	StatusPending  VenueStatus = "pending"
	StatusApproved VenueStatus = "approved"
	StatusRejected VenueStatus = "rejected"
)
