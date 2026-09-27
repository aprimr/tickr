package auth

import (
	"time"

	appMiddleware "github.com/aprimr/tickr/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler AuthHandler) {
	r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/users", handler.HandleUserRegister)
	r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/venues", handler.HandleVenueAdminRegister)

	r.Route("/auth", func(r chi.Router) {
		r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/login", handler.HandleLogin)
		r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/verify-email", handler.HandleVerifyAccount)

		r.With(appMiddleware.RateLimit(5, 5*time.Minute)).Post("/forgot-password", handler.HandleForgotPassword)
	})
}
