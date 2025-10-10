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

	// Check required variables
	requiredVars := []string{"JWT_SECRET"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Required environment variable  %s is not set", v)
		}
	}
}

func main() {
	// Load environment variables
	loadEnv()

	// Set up databases
	var refreshTokenDb = db.SetupRefreshTokenDB()
	var userDb = db.SetupTestUserDB()

	// Create repositories
	tokenRepo := repository.NewRefreshTokenRepository(refreshTokenDb)
	userRepo := repository.NewUserRepository(userDb)

	// Parse token duration
	TTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalf("Error parsing duration:%v", err.Error())
	}

	// Create services
	authService := auth.NewAuthService(userRepo, tokenRepo, os.Getenv("JWT_SECRET"), TTL)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)

	// Create New Router
	r := mux.NewRouter()

	// Public Routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to Golden Market!\n")
	})

	// Public pages
	r.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/login.html")
	}).Methods("GET")

	r.HandleFunc("/api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/register.html")
	}).Methods("GET")

	// Auth API
	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/refresh", authHandler.RefreshToken).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()

	protected.Use(middleware.AuthMiddleware(authService))

	protected.HandleFunc("/profile", userHandler.Profile).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:8080"
	}
	log.Printf("Server starting on %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
