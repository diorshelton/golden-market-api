package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/diorshelton/golden-market/auth"
	"github.com/diorshelton/golden-market/db"
	"github.com/diorshelton/golden-market/handlers"
	"github.com/diorshelton/golden-market/middleware"
	"github.com/diorshelton/golden-market/models"
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
	tokenRepo := models.NewRefreshTokenRepository(refreshTokenDb)
	userRepo := models.NewUserRepository(userDb)

	// Create services
	authService := auth.NewAuthService(userRepo, tokenRepo, os.Getenv("JWT_SECRET"), 15*time.Minute)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)

	// Main serve Mux
	mainMux := http.NewServeMux()

	// Public routes
	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to Golden Market!\n")
	})

	mainMux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/login.html")
	})

	// Register protected endpoints
	mainMux.HandleFunc("POST /register", authHandler.Register)
	mainMux.HandleFunc("POST /login", authHandler.Login)
	mainMux.HandleFunc("POST /refresh", authHandler.RefreshToken)

	mainMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/register.html")
	})

	log.Print("Golden Market server running...")
	log.Fatal(http.ListenAndServe(":3000", mainMux))
}
