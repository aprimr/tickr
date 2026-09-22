package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/aprimr/tickr/internal/auth"
	"github.com/aprimr/tickr/internal/db"
	"github.com/aprimr/tickr/internal/email"
	"github.com/aprimr/tickr/internal/health"
	appMiddleware "github.com/aprimr/tickr/internal/middleware"
	"github.com/aprimr/tickr/internal/worker"
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

	// Init slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Root context for server shutdown signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init worker pool
	workerPool := worker.NewPool(100, logger)
	workerPool.Start(ctx, 4)
	defer workerPool.Stop()

	// Init email service
	_ = email.InitBrevo(logger)
	// workerPool.Enqueue(
	// 	func(ctx context.Context) error {
	// 		body, _ := email.AccountVerification("Hello", "99009")

	// 		return mailer.SendEmail(
	// 			"aprimregmi24@gmail.com",
	// 			"testing",
	// 			body,
	// 		)
	// 	},
	// )

	// Connect to database
	var dbPool *pgxpool.Pool
	if dbPool, err = db.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	r := chi.NewRouter()

	// CORS middleware
	r.Use(appMiddleware.SetupCors(logger))

	// Global middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.ClientIPFromXFF("10.0.0.0/8"))

	// Auth Route Dependencies
	authRepo := auth.NewAuthRepository(dbPool)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService, logger)

	// Health Route Dependency
	healthHandler := health.NewHealthHandler(dbPool)

	// Routes
	r.Route("/v1", func(r chi.Router) {
		auth.RegisterRoutes(r, authHandler)
	})

	// Health route
	health.RegisterRoutes(r, healthHandler)

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
