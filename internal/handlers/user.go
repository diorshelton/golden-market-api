package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/diorshelton/golden-market-api/internal/middleware"
	"github.com/diorshelton/golden-market-api/internal/repository"
)

// UserHandler contains HTTP handlers for user-related endpoints
type UserHandler struct {
	userRepo *repository.UserRepository
}

// NewUserHandler creates a new yser handler
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// UserResponse represents the user data returned to clients
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	FirstName	string `json:"first_name"`
	LastName string `json:"last_name"`
	Balance  int64 `json:"balance"`
	Inventory []string `json:"inventory"`
	CreatedAt  string `json:"created_at"`
}

// Profile returns the authenticated user's profile
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request context (set by auth middleware)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user, err := h.userRepo.GetUserProfile(userID)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	//TODO:Add inventory to user profile
	// Return user profile data
	response := UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email:    user.Email,
		Balance:  int64(user.Balance),
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
