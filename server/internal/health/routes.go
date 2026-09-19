package health

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, handler HealthHandler) {
	r.Get("/health", handler.HandleHealth)
}
