package auth

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, handler AuthHandler) {
	r.Post("/users", handler.HandleUserRegister)
	r.Post("/venues", handler.HandleVenueAdminRegister)
}
