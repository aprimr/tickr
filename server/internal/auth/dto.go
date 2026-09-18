package auth

import "github.com/aprimr/tickr/internal/pkg/validate"

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterRequest struct {
	FullName    string `json:"fullname" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type VenueRegisterRequest struct {
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`

	VenueName    string `json:"venue_name" validate:"required"`
	Address      string `json:"address" validate:"required"`
	City         string `json:"city" validate:"required"`
	TotalScreens int    `json:"total_screens" validate:"required, min=1"`
}

// Validate UserRegisterRequest data
func (req *UserRegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if !validate.IsMinLength(req.FullName, 4) {
		errs["fullname"] = "full name must be greater than 4 characters"
	}
	if !validate.IsEmail(req.Email) {
		errs["email"] = "invalid email format"
	}
	if !validate.IsPhone(req.PhoneNumber) {
		errs["phone_number"] = "invalid phone number format"
	}
	if !validate.IsValidPassword(req.Password) {
		errs["password"] = "password must be at least 8 characters long and contain at least one number and one special character (! @ # $ % ^ & * ? )"
	}

	return errs
}

// Validate VenueRegisterRequest data
func (req *VenueRegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if !validate.IsEmail(req.Email) {
		errs["email"] = "invalid email format"
	}
	if !validate.IsPhone(req.PhoneNumber) {
		errs["phone_number"] = "invalid phone number format"
	}
	if !validate.IsValidPassword(req.Password) {
		errs["password"] = "password must be at least 8 characters long and contain at least one number and one special character (! @ # $ % ^ & * ? )"
	}
	if !validate.IsMinLength(req.VenueName, 5) {
		errs["venue_name"] = "venue name must be atleast 5 characters long"
	}
	if !validate.IsNotBlank(req.Address) {
		errs["address"] = "address is required"
	}
	if !validate.IsNotBlank(req.City) {
		errs["city"] = "city is required"
	}
	if req.TotalScreens < 1 {
		errs["total_screens"] = "total screens must be at least 1"
	}

	return errs
}
