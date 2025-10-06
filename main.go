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
	tokenRepo := models.NewRefreshTokenRepository(refreshTokenDb)
	userRepo := models.NewUserRepository(userDb)

	// Create services
	authService := auth.NewAuthService(userRepo, tokenRepo, os.Getenv("JWT_SECRET"), 15*time.Minute)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)

	// Create New Router
	r := mux.NewRouter()

	// Public Routes
	r.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to Golden Market!\n")
	})

	// Public pages
	r.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/login.html")
	}).Methods("GET")

	r.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/register.html")
	}).Methods("GET")

	// Auth API
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/refresh", authHandler.RefreshToken).Methods("POST")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	protected.HandleFunc("GET /profile", userHandler.Profile).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func sanityCheck(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println("Form Data:", r.Form)
	log.Println("PostForm Data:", r.PostForm)
}
