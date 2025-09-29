package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// CreateRefreshToken creates a new refresh token for a user
func (r *RefreshTokenRepository) CreateRefreshToken(userID uuid.UUID, ttl time.Duration) (*RefreshToken, error) {
	// Generate a unique token identifier
	tokenID := uuid.New()
	expiresAt := time.Now().Add(ttl)

	token := &RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     tokenID.String(), // User the UUID as the token
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	query := `
INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked) VALUES (?,?,?,?,?,?)
`
	_, err := r.db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt, token.Revoked)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetRefreshToken retrieves a refresh token by its token string
func (r *RefreshTokenRepository) GetRefreshToken(tokenString string) (*RefreshToken, error) {
	query := `
	SELECT id, user_id, token, expires_at, created_at, revoked
	FROM refresh_tokens
	WHERE token = ?
	`

	var token RefreshToken
	err := r.db.QueryRow(query, tokenString).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.Revoked,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (r *RefreshTokenRepository) RevokeRefreshToken(tokenString string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE token = ?
	`

	_, err := r.db.Exec(query, tokenString)
	return err
}
