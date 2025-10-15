package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/diorshelton/golden-market/internal/auth"
	"github.com/diorshelton/golden-market/internal/database"
	"github.com/diorshelton/golden-market/internal/handlers"
	"github.com/diorshelton/golden-market/internal/middleware"
	"github.com/diorshelton/golden-market/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// loadEnv loads environment variables from .env file
func loadEnv() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, using environment variable")
	}

	// Check all required variables
	requiredVars := []string{
		"JWT_SECRET",
		"REFRESH_SECRET",
		"ACCESS_TOKEN_EXPIRY",
	}

	missing := []string{}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Required environment variable  %s is missing", missing)
	}
}

func main() {
	// Load environment variables
	loadEnv()

	// Set up databases
	refreshTokenDb := db.SetupRefreshTokenDB()
	userDb := db.SetupTestUserDB()

	// Create repositories
	tokenRepo := repository.NewRefreshTokenRepository(refreshTokenDb)
	userRepo := repository.NewUserRepository(userDb)

	// Parse token duration
	accessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalf("Invalid ACCESS_TOKEN_EXPIRY: %v", err)
	}

	refreshTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalf("Invalid REFRESH_TOKEN_EXPIRY: %v", err)
	}

	// Create  auth service
	authService := auth.NewAuthService(
		userRepo,
		tokenRepo,
		os.Getenv("JWT_SECRET"),
		os.Getenv("REFRESH_SECRET"),
		accessTTL,
		refreshTTL,
	)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)

	// Create router
	r := mux.NewRouter()

	// Public Routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to Golden Market!\n")
	})

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	}).Methods("GET")

	// Public pages
	r.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/login.html")
	}).Methods("GET")

	r.HandleFunc("/api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/register.html")
	}).Methods("GET")

	// Auth API Endpoints
	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/refresh", authHandler.Refresh).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))
	protected.HandleFunc("/profile", userHandler.Profile).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
