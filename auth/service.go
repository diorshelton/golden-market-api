package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/diorshelton/golden-market/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
  ErrInvalidCredentials = errors.New("invalid credentials")
  ErrInvalidToken       = errors.New("invalid token")
  ErrExpiredToken       = errors.New("token has expired")
  ErrEmailInUse         = errors.New("email already in use")
)

//AuthService provides authentication functionality
type AuthService struct {
	userRepo *models.UserRepository
	refreshTokenRepo *models.RefreshTokenRepository
	jwtSecret []byte
	accessTokenTTL time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *models.UserRepository, refreshTokenRepo *models.RefreshTokenRepository, jwtSecret string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret: []byte(jwtSecret),
		accessTokenTTL: accessTokenTTL,
	}
}

//Register creates a new user with the provided credentials
func (s *AuthService) Register(firstName, lastName, email, username, password string)(*models.User, error) {
	//Check if user already exists
	_,err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return  nil, ErrEmailInUse
	}

	//Only proceed if the error was "user not found"
	if !errors.Is(err, sql.ErrNoRows) {
		return  nil, err
	}

	//Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	//Create the user
	user, err := s.userRepo.CreateUser(username, firstName , lastName, email, password)
	if err != nil {
		return  nil,err
	}
	return  user, nil
}

func (s *AuthService) Login (email, password string) (string, error) {
	// Get the user from the database
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Verify the password
	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return  " ", ErrInvalidCredentials
	}

	// Generate an access token
	token, err := s.generateAccessToken(user)
	if err != nil {
		return  "", err
	}

	return  token, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	// Set the expiration time
	expirationTime := time.Now().Add(s.accessTokenTTL)

	// Create the JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"username": user.Username,
		"email": user.Email,
		"exp": expirationTime.UTC(),
		"iat": time.Now().UTC(),
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Sign the token with secret key
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken verifies a JWT token and returns the claim
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// Validate the signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return  nil, ErrInvalidToken
	}
	return  s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return  nil, ErrExpiredToken
		}
		return  nil, ErrInvalidToken
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}