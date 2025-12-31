package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/diorshelton/golden-market-api/internal/auth"
	"github.com/diorshelton/golden-market-api/internal/models"
)

type AuthServiceInterface interface {
	Register(firstName, lastName, email, username, password string) (*models.User, error)
	Login(email, password string) (string, string, error)
	Refresh(oldRefreshToken string) (*auth.TokenPair, error)
	Logout(refreshToken string) error
}

// AuthHandler contains HTTP handlers for authentication
type AuthHandler struct {
	authService AuthServiceInterface
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Username        string `json:"username"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// Validates form input from user's POST request
func (r *RegisterRequest) Validate() error {
	// Trim whitespace from all fields
	r.Username = strings.TrimSpace(r.Username)
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.PasswordConfirm = strings.TrimSpace(r.PasswordConfirm)

	// Check values for empty strings
	if r.Username == "" || r.FirstName == "" || r.LastName == "" || r.Email == "" || r.Password == "" {
		return errors.New("all fields required")
	}

	if len(r.Username) < 3 || len(r.Username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}

	if len(r.Password) < 8 || len(r.Password) > 64 {
		return errors.New("password must be between 8 and 64 characters")
	}

	// Verify both passwords match
	if r.Password != r.PasswordConfirm {
		return errors.New("passwords must match")
	}

	// Validate email address
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

// RegisterResponse contains the user data after successful registration
type RegisterResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest
	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the auth service
	user, err := h.authService.Register(req.FirstName, req.LastName, req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailInUse) {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating user", http.StatusConflict)
		log.Printf("Error occurred: %v", err)
		return
	}

	// Return response
	response := RegisterResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Message:   "Registration successful",
	}

	log.Printf("Registering:%v", response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	// Trim whitespace
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	// Check for empty credentials
	if r.Email == "" || r.Password == "" {
		return errors.New("invalid credentials")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}

// LoginResponse contains the JWT token after successful login
type LoginResponse struct {
	AccessToken string `json:"token"`
}

// Login handles user login with access and refresh tokens
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse JSON
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate and normalize
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Attempt to login
	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	maxAge := int((7 * 24 * time.Hour).Seconds())
	// Set refresh token in HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
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
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	maxAge := int((7 * 24 * time.Hour).Seconds())
	// Set new refresh token in HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken, //New rotated token
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	// Return the new access token
	response := RefreshResponse{Token: tokenPair.AccessToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout of current device
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Read refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // Expire the cookie
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
		return
	}

	refreshToken := cookie.Value

	// Delete the refresh token from the database
	err = h.authService.Logout(refreshToken)
	if err != nil {
		log.Printf("Failed to delete refresh. token: %v", err)
		// Continue to clear cookie even if error occurs
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire cookie
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
