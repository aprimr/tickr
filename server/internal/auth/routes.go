package auth

import (
	"time"

	appMiddleware "github.com/aprimr/tickr/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler AuthHandler) {
	// Global limit, 300 requests combined on all routes in a minute
	r.Use(appMiddleware.RateLimit(300, time.Minute))

	r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/users", handler.HandleUserRegister)
	r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/venues", handler.HandleVenueAdminRegister)

	r.Route("/auth", func(r chi.Router) {
		r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/login", handler.HandleLogin)
		r.With(appMiddleware.RateLimit(10, 5*time.Minute)).Post("/verify-email", handler.HandleVerifyAccount)

		r.With(appMiddleware.RateLimit(5, 15*time.Minute)).Post("/forgot-password", handler.HandleForgotPassword)
		r.With(appMiddleware.RateLimit(5, 15*time.Minute)).Post("/reset-password", handler.HandleResetPassword)

		r.With(appMiddleware.RateLimit(2, time.Minute)).Post("/rotate", handler.HandleTokenRotation)

		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Authenticate()) // Allow logined users with any role

			r.With(appMiddleware.RateLimit(10, time.Minute)).Post("/logout", handler.HandleLogout)
			r.With(appMiddleware.RateLimit(10, time.Minute)).Post("/logout-all", handler.HandleLogoutAllDevices)
		})
	})
}
