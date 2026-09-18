package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aprimr/tickr/internal/auth"
	"github.com/aprimr/tickr/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Not found or failed to load env: %v\n", err)
	}

	// Connect to database
	var dbPool *pgxpool.Pool
	if dbPool, err = db.Connect(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	r := chi.NewRouter()

	// Global middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Auth Dependencies
	authRepo := auth.NewAuthRepository(dbPool)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Routes
	r.Route("/v1", func(r chi.Router) {
		auth.Routes(r, authHandler)
	})

	// Health route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := dbPool.Ping(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"error", "database":"unreachable"}`))

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok", "database":"connected"}`))
	})

	// Start the server and listen to port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s...", port)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
