package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/diorshelton/golden-market/auth"
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
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
}

// RegisterResponse contains the user data after successful registration
type RegisterResponse struct {
	ID string `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
// parse the request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

// Parse registration form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "An error occurred while parsing form: %v", http.StatusInternalServerError)
	}

// Get form input and trim white space if any
	username := strings.TrimSpace(r.Form.Get("username"))
	firstName := strings.TrimSpace(r.Form.Get("first_name"))
	lastName := strings.TrimSpace(r.Form.Get("last_name"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))
	passwordConfirm := strings.TrimSpace(r.Form.Get("password_confirm"))

	// Validate input
	err := validateInput(username, firstName, lastName, email, password,passwordConfirm, w)
	if err != nil {
		http.Error(w, "A parsing/validation error occurred", http.StatusBadRequest)
	}

// Call the auth service to register the user
	user, err := h.authService.Register(firstName, lastName, email, username,password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailInUse) {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Return created user (without sensitive data)
	response := RegisterResponse {
		ID: user.ID.String(),
		Email: user.Email,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Validates form input from user's POST request
func validateInput(username,firstName, lastName, email, password, passwordConfirm string, w http.ResponseWriter) error {
	// Check values for empty strings
	if username == "" || firstName == "" || lastName == "" || email == "" || password == ""{
		http.Error(w, "Required form input empty", http.StatusBadRequest)
		return errors.New("missing required input")
	}

	// Verify both passwords match
	if password != passwordConfirm {
		http.Error(w, "passwords must match", http.StatusBadRequest)
		return errors.New("passwords must match")
	}

	// Validate email address
	if _, err := mail.ParseAddress(email); err != nil {
		http.Error(w, "invalid email address", http.StatusBadRequest)
		return errors.New("invalid email address")
	}
	return  nil
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse contains the JWT token after successful login
type LoginResponse struct {
	Token string `json:"toke"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w,"Invalid request payload", http.StatusBadRequest)
	}
	// Parse login input
	r.ParseForm()
		if err := r.ParseForm(); err != nil {
		http.Error(w, "An error occurred while parsing form: %v", http.StatusInternalServerError)
	}

	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))

	// Attempt to login
	token, err := h.authService.Login(email, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the token
	response := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshRequest represents the refresh token payload
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshResponse contains the new access token
type RefreshResponse struct {
	Token string `json:"token"`
}

// RefreshToken handles access token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Attempt to refresh the token
	token, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken){
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return the new access token
	response := RefreshResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}