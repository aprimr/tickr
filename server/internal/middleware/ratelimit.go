package middleware

import (
	"net/http"
	"time"

	"github.com/aprimr/tickr/internal/pkg/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

// RateLimit returns a preconfigured rate limiter middleware
func RateLimit(reqLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return httprate.LimitBy(
		reqLimit,
		windowLength,
		func(r *http.Request) (string, error) {
			return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
		},
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			response.Error(w, http.StatusTooManyRequests, "too many requests, please try again later", nil)
		}),
	)
}
