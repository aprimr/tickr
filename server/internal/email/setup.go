package email

import (
	"context"
	"log/slog"
	"os"

	"github.com/aprimr/tickr/internal/worker"
	brevo "github.com/getbrevo/brevo-go/lib"
)

type EmailService interface {
	SendEmail(to string, subject string, body string) error
}

type emailService struct {
	client    *brevo.APIClient
	fromEmail string
	fromName  string
	pool      *worker.Pool
	logger    *slog.Logger
}

// NewEmailService initializes and returns the EmailService interface
func NewEmailService(pool *worker.Pool, logger *slog.Logger) EmailService {
	// Get envs
	apiKey := os.Getenv("BREVO_API_KEY")
	fromEmail := os.Getenv("BREVO_FROM_EMAIL")
	fromName := os.Getenv("BREVO_FROM_NAME")

	// Configure brevo
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	client := brevo.NewAPIClient(cfg)

	return &emailService{
		client:    client,
		fromEmail: fromEmail,
		fromName:  fromName,
		pool:      pool,
		logger:    logger,
	}
}

func (s *emailService) SendEmail(to string, subject string, body string) error {
	// Adds the email job to the pool
	s.pool.Enqueue(func(ctx context.Context) error {
		apiCtx := context.Background()

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

		_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(apiCtx, emailData)
		if err != nil {
			s.logger.Error("failed to send email", "error", err, "email subject", subject, "to", to)
			return err
		}

		return nil
	})

	return nil
}
