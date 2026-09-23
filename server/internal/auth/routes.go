package auth

import (
	"time"

	appMiddleware "github.com/aprimr/tickr/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler AuthHandler) {
	r.Use(appMiddleware.RateLimit(10, time.Minute))

	r.Post("/users", handler.HandleUserRegister)
	r.Post("/venues", handler.HandleVenueAdminRegister)

	r.Post("/verify", handler.HandleVerifyAccount)
}
