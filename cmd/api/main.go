package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/diorshelton/golden-market-api/internal/auth"
	"github.com/diorshelton/golden-market-api/internal/database"
	"github.com/diorshelton/golden-market-api/internal/handlers"
	"github.com/diorshelton/golden-market-api/internal/middleware"
	"github.com/diorshelton/golden-market-api/internal/repository"
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
	database := database.SetupTestDB()
	defer database.Close()

	// Create repositories
	tokenRepo := repository.NewRefreshTokenRepository(database)
	userRepo := repository.NewUserRepository(database)

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

	//Apply CORS middleware
	r.Use(middleware.CORS)

	// Public Routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to Golden Market!\n")
	})

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":      "ok",
			"port":        os.Getenv("PORT"),
			"environment": os.Getenv("ENVIRONMENT"),
		})
	}).Methods("GET")

	// Public pages
	r.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/login.html")
	}).Methods("GET")

	r.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/register.html")
	}).Methods("GET")

	// --- Auth API Endpoints (rate limited) ---
	authRouter := r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.Use(middleware.CORS)      // Apply CORS to Subrouter
	authRouter.Use(middleware.RateLimit) // Apply ratelimiting

	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/refresh", authHandler.Refresh).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")

	// --- Protected routes ---
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.CORS) // Apply CORS to Subrouter
	protected.Use(middleware.Auth(authService))
	protected.HandleFunc("/profile", userHandler.Profile).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s", os.Getenv("ENVIRONMENT"))

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
