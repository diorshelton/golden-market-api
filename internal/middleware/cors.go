package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS handles Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		var allowedOrigins []string

		if allowedOriginsEnv != "" {
			allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		} else {
			// Default allowed origins for development
			allowedOrigins = []string{
				"http://localhost:5173", // Vite default
				"http://localhost:3000", // React default
				"http://localhost:8080",
			}
		}

		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.TrimSpace(allowedOrigin) == origin {
				allowed = true
				break
			}
		}

		// Set CORS headers if origin is allowed
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}