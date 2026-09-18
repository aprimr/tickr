package auth

import "github.com/go-chi/chi/v5"

func Routes(r chi.Router, handler AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/users/register", handler.HandleUserRegister)
		r.Post("/venues/register", handler.HandleVenueAdminRegister)
	})
}
