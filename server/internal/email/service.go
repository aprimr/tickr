package email

import (
	"bytes"
	"fmt"

	"github.com/aprimr/tickr/internal/email/templates"
)

// AccountVerificationEmail builds the email body and sends it
func (s *emailService) SendAccountVerificationEmail(to, name, otp string) error {
	data := VerificationData{
		Name: name,
		OTP:  otp,
	}

	var buf bytes.Buffer
	if err := templates.AccountVerification.Execute(&buf, data); err != nil {
		s.logger.Error("failed to send account verification email", "error", err)
		return fmt.Errorf("failed to build email body: %w", err)
	}

	emailSubject := "Verify Your Tickr Account"
	emailBody := buf.String()

	return s.SendEmail(to, emailSubject, emailBody)
}
