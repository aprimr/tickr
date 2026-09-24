package validate

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^9[0-9]{9}$`)
)

// IsEmail checks if a string is valid email
func IsEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsPhone checks if a string is valid phone number
//
// Must be 10 digits and start with 9
func IsPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// IsValidPassword checks if given string follows the password requirements
//
// Valid Password must:
//
//   - be greater or equal to 8 characters and less or equal to 40 characters,
//   - contain atleast one number,
//   - contain atleast one special character (! @ # $ % ^ & * ? )
func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 40 {
		return false
	}
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*?]`).MatchString(password)
	return hasNumber && hasSpecial
}

// IsMinLength checks if a trimmed string is strictly greater than the minimum length
func IsMinLength(s string, min int) bool {
	return len(strings.TrimSpace(s)) > min
}

// IsNotBlank checks if a string is not empty or just whitespace
func IsNotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}
