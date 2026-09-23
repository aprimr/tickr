package auth

import (
	"github.com/aprimr/tickr/internal/pkg/validate"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type VenueRegisterRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`

	VenueName    string `json:"venue_name"`
	Address      string `json:"address"`
	City         string `json:"city"`
	TotalScreens int    `json:"total_screens"`
}

type VerifyAccountRequest struct {
	UserID uuid.UUID `json:"user_id"`
	OTP    string    `json:"otp"`
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

// Validate VerifyAccountRequest data
func (req *VerifyAccountRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if req.UserID == uuid.Nil {
		errs["user_id"] = "user_id is required"
	}
	if req.OTP == "" {
		errs["otp"] = "otp is required"
	} else if len(req.OTP) != 6 {
		errs["otp"] = "otp must be 6 characters long"
	}
	return errs
}
