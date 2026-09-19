package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

// SetupCors initializes CORS and returns the handler
func SetupCors(logger *slog.Logger) func(next http.Handler) http.Handler {

	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if originsEnv != "" {
		allowedOrigins = strings.Split(originsEnv, ",")
	} else {
		logger.Error("ALLOWED_ORIGINS environment variable is empty")
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
