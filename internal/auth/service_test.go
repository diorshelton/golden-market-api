package auth

import (
	"testing"
	"time"

	"github.com/diorshelton/golden-market/internal/database"
	// "github.com/diorshelton/golden-market/internal/models"
	"github.com/diorshelton/golden-market/internal/repository"
	// "github.com/google/uuid"
)

func setupTestService(t *testing.T) (*AuthService, func()) {
	db := database.SetupTestDB()

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)

	service := NewAuthService(
		userRepo,
		tokenRepo,
		"test_jwt_secret",
		"test_refresh_secret",
		time.Minute*15,
		time.Hour*24*7,
	)

	cleanup := func() {
		database.CleanupTestDB(db)
	}

	return service, cleanup
}

func TestRegister(t * testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	tests := []struct {
		name string
		firstName string
		lastName string
		email string
		username string
		password string
		wantErr error
	}{
		{
			name: "successful registration",
			firstName: "Jake",
			lastName: "The Dog",
			email: "jdog@example.com",
			username: "jdog",
			password: "password123",
			wantErr:  nil,
		},
		{
			name: "duplicate username",
			firstName: "Joshua",
			lastName: "The Dog",
			email: "jake@example.com",
			username: "jdog",
			password: "password",
			wantErr:  ErrUsernameExists,
		},
		{
			name: "duplicate email",
			firstName: "Joshua",
			lastName: "The Dog",
			email: "jdog@example.com",
			username: "jodogo",
			password: "password123",
			wantErr:  ErrEmailInUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t * testing.T) {
			user, err := service.Register(tt.firstName, tt.lastName, tt.email, tt.username, tt.password)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				} else if err != tt.wantErr {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if user.Username != tt.username{
				t.Errorf("Expected username %s, got %s", tt.username,
				user.Username)
			}
			if user.Email != tt.email {
				t.Errorf("Expected email %s, got %s", tt.email, user.Email)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	// Create a user to test login
	testPassword := "adventuretime"
	user, err := service.Register("Finn", "The Human", "ooofinn@example.com", "finn", "adventuretime")
	if err != nil {
		t.Fatalf("Failed to register user for login test: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "successful login",
			email:    user.Email,
			password: testPassword,
			wantErr:  nil,
		},
		{
			name: "wrong password",
			email: user.Email,
			password: "wrongpassword",
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "non-existent email",
			email: "nonexistant@example.com",
			password: "password",
			wantErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken, refreshToken, err := service.Login(tt.email, tt.password)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if accessToken == "" {
				t.Error("Expected access token, got empty string")
			}
			if refreshToken == "" {
				t.Error("Expected refresh token, got empty string")
			}

			//Verify token was stored in DB
			tokenRecord, err := service.refreshTokenRepo.GetRefreshToken(refreshToken)
			if err != nil {
				t.Fatalf("Failed to retrieve refresh token from DB: %v", err)
			}
			if tokenRecord.UserID != user.ID {
				t.Errorf("Expected token UserID %v, got %v", user.ID, tokenRecord.UserID)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	// Create test user and login
	user, _ := service.Register("Cosmic", "Owl", "cosmico@example.com", "cosmico", "password123")
	_, refreshToken, _ := service.Login(user.Email, "password123")

	tests := []struct {
		name     string
		token    string
		wantErr  error
	}{
		{
			name:    "successful refresh",
			token:   refreshToken,
			wantErr: nil,
		},
		{
			name:    "invalid token",
			token:   "invalid_token",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenPair, err := service.Refresh(tt.token)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tokenPair.AccessToken == "" {
				t.Error("Expected new access token")
			}
			if tokenPair.RefreshToken == "" {
				t.Error("Expected new refresh token")
			}

			// Verify old token was deleted
			_, err = service.refreshTokenRepo.GetRefreshToken(tt.token)
			if err == nil {
				t.Error("Expected old token to be deleted")
			}

			// Verify new token exists
			newTokenRecord, err := service.refreshTokenRepo.GetRefreshToken(tokenPair.RefreshToken)
			if err != nil {
				t.Fatalf("Failed to retrieve new refresh token: %v", err)
			}
			if newTokenRecord.UserID != user.ID {
				t.Error("New token should belong to same user")
			}
		})
	}
}

func TestRefreshWithExpiredToken(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	// Create test user
	user, _ := service.Register("Marceline", "Abadeer", "MarcelineTheVampireQueen@example.com", "marcelinequeen", "password123")

	// Create an expired refresh token
	expiredToken, _ := service.refreshTokenRepo.CreateRefreshToken(user.ID, -1*time.Hour)

	_, err := service.Refresh(expiredToken.Token)
	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got %v", err)
	}
}

func TestLogout(t *testing.T) {
	service, cleanup := setupTestService(t)
	defer cleanup()

	// Create test user and login
	user, _ := service.Register("Peppermint", "Butler", "Pepbut@example.com", "pepbut", "password123")
	_, refreshToken, _ := service.Login(user.Email, "password123")

	// Logout
	err := service.Logout(refreshToken)
	if err != nil {
		t.Fatalf("Unexpected error during logout: %v", err)
	}

	// Verify token was deleted
	_, err = service.refreshTokenRepo.GetRefreshToken(refreshToken)
	if err == nil {
		t.Error("Expected token to be deleted after logout")
	}
}