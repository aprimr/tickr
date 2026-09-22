package email

import (
	"context"
	"log/slog"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type EmailService interface {
	SendEmail(to string, subject string, body string) error
}

type BrevoService struct {
	client    *brevo.APIClient
	fromEmail string
	fromName  string
	logger    *slog.Logger
}

func InitBrevo(logger *slog.Logger) *BrevoService {
	// Get envs
	apiKey := os.Getenv("BREVO_API_KEY")
	fromEmail := os.Getenv("BREVO_FROM_EMAIL")
	fromName := os.Getenv("BREVO_FROM_NAME")

	// Configure brevo
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	client := brevo.NewAPIClient(cfg)

	return &BrevoService{
		client:    client,
		fromEmail: fromEmail,
		fromName:  fromName,
		logger:    logger,
	}
}

func (s *BrevoService) SendEmail(to string, subject string, body string) error {
	ctx := context.Background()

	emailData := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Email: s.fromEmail,
			Name:  s.fromName,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: to,
			},
		},
		Subject:     subject,
		HtmlContent: body,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailData)
	if err != nil {
		s.logger.Error("failed to send email", "error", err, "email subject", subject, "to", to)
		return err
	}

	return nil
}
