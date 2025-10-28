package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/diorshelton/golden-market/internal/auth"
	"github.com/diorshelton/golden-market/internal/models"
	"github.com/google/uuid"
)

// Mock AuthService for testing
type MockAuthService struct {
	RegisterFunc func(firstName, lastName, email, username, password string) (*models.User, error)
	LoginFunc    func(email, password string) (accessToken string, refreshToken string, err error)
	RefreshFunc  func(oldRefreshToken string) (*auth.TokenPair, error)
	LogoutFunc   func(tokenString string) error
}

func (m *MockAuthService) Register(firstName, lastName, email, username, password string) (*models.User, error) {
	return m.RegisterFunc(firstName, lastName, email, username, password)
}

func (m *MockAuthService) Login(email, password string) (string, string, error) {
	return m.LoginFunc(email, password)
}

func (m *MockAuthService) Refresh(oldRefreshToken string) (*auth.TokenPair, error) {
	return m.RefreshFunc(oldRefreshToken)
}

func (m *MockAuthService) Logout(tokenString string) error {
	return m.LogoutFunc(tokenString)
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		mockRegister   func(firstName, lastName, email, username, password string) (*models.User, error)
		expectedStatus int
		checkResponse  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "successful registration",
			formData: url.Values{
				"first_name":       []string{"Dandara"},
				"last_name":        []string{"dos Palmares"},
				"email":            []string{"dandap@example.com"},
				"username":         []string{"dandap"},
				"password":         []string{"password123"},
				"password_confirm": []string{"password123"},
			},
			mockRegister: func(firstName, lastName, email, username, password string) (*models.User, error) {
				return &models.User{
					ID:        uuid.New(),
					FirstName: firstName,
					LastName:  lastName,
					Email:     email,
					Username:  username,
				}, nil
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var response RegisterResponse
				err := json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response.Message != "Registration successful" {
					t.Errorf("Expected message 'Registration successful', got %s", response.Message)
				}
			},
		},
		{
			name: "missing required fields",
			formData: url.Values{
				"first_name": []string{"Dandara"},
				"email":      []string{"dandap@example.com"},
			},
			mockRegister:   nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "passwords don't match",
			formData: url.Values{
				"first_name":       []string{"Dandara"},
				"last_name":        []string{"dos Palmares"},
				"email":            []string{"dandap@example.com"},
				"username":         []string{"dandap"},
				"password":         []string{"password123"},
				"password_confirm": []string{"different"},
			},
			mockRegister:   nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "duplicate username",
			formData: url.Values{
				"first_name":       []string{"Dandara"},
				"last_name":        []string{"dos Palmares"},
				"email":            []string{"dandap@example.com"},
				"username":         []string{"dandap"},
				"password":         []string{"password123"},
				"password_confirm": []string{"password123"},
			},
			mockRegister: func(firstName, lastName, email, username, password string) (*models.User, error) {
				return nil, auth.ErrUsernameExists
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				RegisterFunc: tt.mockRegister,
			}
			handler := NewAuthHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/auth/register", nil)
			req.Form = tt.formData
			req.PostForm = tt.formData

			rr := httptest.NewRecorder()
			handler.Register(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		mockLogin      func(email, password string) (string, string, error)
		expectedStatus int
		checkCookie    bool
	}{
		{
			name: "successful login",
			formData: url.Values{
				"email":    []string{"john@example.com"},
				"password": []string{"password123"},
			},
			mockLogin: func(email string, password string) (string, string, error) {
				return "access_token_here",
					"refresh_token_here", nil
			},
			expectedStatus: http.StatusOK,
			checkCookie:    true,
		},
		{
			name: "invalid credentials",
			formData: url.Values{
				"email":    []string{"john@example.com"},
				"password": []string{"wrongpassword"},
			},
			mockLogin: func(email string, password string) (string, string, error) {
				return "", "", auth.ErrInvalidCredentials
			},
			expectedStatus: http.StatusUnauthorized,
			checkCookie:    false,
		},
		{
			name: "missing email",
			formData: url.Values{
				"email":    []string{" "},
				"password": []string{" "},
			},
			mockLogin:      nil,
			expectedStatus: http.StatusUnauthorized,
			checkCookie:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				LoginFunc: tt.mockLogin,
			}
			handler := NewAuthHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
			req.Form = tt.formData
			req.PostForm = tt.formData

			rr := httptest.NewRecorder()
			handler.Login(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.checkCookie {
				cookies := rr.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						found = true
						if !cookie.HttpOnly {
							t.Error("Expected refresh_token cookie to be HttpOnly")
						}
						if cookie.SameSite != http.SameSiteStrictMode {
							t.Error("Expected refresh_token cookie to have SameSite=Strict")
						}
					}
				}
				if !found {
					t.Error("Expected refresh_token cookie to be set")
				}
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(req *http.Request)
		mockRefresh    func(oldRefreshToken string) (*auth.TokenPair, error)
		expectedStatus int
	}{
		{
			name: "successful token refresh",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "valid_refresh_token",
				})
			},
			mockRefresh: func(oldRefreshToken string) (*auth.TokenPair, error) {
				return &auth.TokenPair{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing refresh token cookie",
			setupCookie:    func(req *http.Request) {},
			mockRefresh:    nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "expired refresh token",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "expired_token",
				})
			},
			mockRefresh: func(oldRereshToken string) (*auth.TokenPair, error) {
				return nil, auth.ErrExpiredToken
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid refresh token",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "invalid_token",
				})
			},
			mockRefresh: func(oldRefreshToken string) (*auth.TokenPair, error) {
				return nil, auth.ErrInvalidToken
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				RefreshFunc: tt.mockRefresh,
			}
			handler := NewAuthHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
			tt.setupCookie(req)

			rr := httptest.NewRecorder()
			handler.Refresh(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(req *http.Request)
		mockLogout     func(tokenString string) error
		expectedStatus int
		checkCookie    bool
	}{
		{
			name: "successful logout",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: "valid_token",
				})
			},
			mockLogout: func(tokenString string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			checkCookie:    true,
		}, {
			name:        "logout without cookie (already logged out)",
			setupCookie: func(req *http.Request) {},
			mockLogout: func(tokenString string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			checkCookie:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				LogoutFunc: tt.mockLogout,
			}
			handler := NewAuthHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			tt.setupCookie(req)

			rr := httptest.NewRecorder()
			handler.Logout(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.checkCookie {
				cookies := rr.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						found = true
						if cookie.MaxAge != -1 {
							t.Errorf("Expected MaxAge -1, got %d", cookie.MaxAge)
						}
					}
				}
				if !found {
					t.Error("Expected refresh_token cookie to be cleared")
				}
			}
		})
	}
}
