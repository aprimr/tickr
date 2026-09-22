package email

import (
	"bytes"

	"github.com/aprimr/tickr/internal/email/templates"
)

// AccountVerificationEmail returns the HTML for account verification
func AccountVerification(name string, otp string) (string, error) {
	data := VerificationData{
		Name: name,
		OTP:  otp,
	}

	var buf bytes.Buffer
	if err := templates.AccountVerification.Execute(&buf, data); err != nil {
		return "", err
	}

	htmlBody := buf.String()

	return htmlBody, nil
}
