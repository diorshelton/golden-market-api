package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/diorshelton/golden-market/internal/auth"
)

// AuthHandler contains HTTP handlers for authentication
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// RegisterResponse contains the user data after successful registration
type RegisterResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse  form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract form data and tri
	username := strings.TrimSpace(r.Form.Get("username"))
	firstName := strings.TrimSpace(r.Form.Get("first_name"))
	lastName := strings.TrimSpace(r.Form.Get("last_name"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))
	passwordConfirm := strings.TrimSpace(r.Form.Get("password_confirm"))

	// Validate input
	if err := validateInput(username, firstName, lastName, email, password, passwordConfirm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the auth service
	user, err := h.authService.Register(firstName, lastName, email, username, password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailInUse) {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Return response
	response := RegisterResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// LoginResponse contains the JWT token after successful login
type LoginResponse struct {
	AccessToken string `json:"token"`
}

// Login handles user login with access and refresh tokens
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))

	// Attempt to login
	accessToken, refreshToken, err := h.authService.Login(email, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Set refresh token in HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	})

	response := LoginResponse{AccessToken: accessToken}

	// Return access tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type RefreshResponse struct {
	Token string `json:"token"`
}

// RefreshToken handles access token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Read old refresh token from  cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}

	oldRefreshToken := cookie.Value

	// Refresh and rotate tokens
	tokenPair, err := h.authService.Refresh(oldRefreshToken)
	if err != nil {
		log.Printf("RefreshAccessToken failed: %v", err)
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		} else {
			log.Printf("Internal server error during token refresh: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	// Set new refresh token in HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken, //New rotated token
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	})
	// Return the new access token
	response := RefreshResponse{Token: tokenPair.AccessToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Validates form input from user's POST request
func validateInput(username, firstName, lastName, email, password, passwordConfirm string) error {
	// Check values for empty strings
	if username == "" || firstName == "" || lastName == "" || email == "" || password == "" {
		return errors.New("all fields required")
	}

	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}

	if len(password) < 8 || len(password) > 16 {
		return errors.New("password must be between 8 and 16 characters")
	}

	// Verify both passwords match
	if password != passwordConfirm {
		return errors.New("passwords must match")
	}

	// Validate email address
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email address")
	}
	return nil
}
