package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/diorshelton/golden-market/internal/models"
	"github.com/google/uuid"
)

// RefreshTokenRepository handles database operations for refresh tokens
type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// CreateRefreshToken creates a new refresh token for a user with secure random token
func (r *RefreshTokenRepository) CreateRefreshToken(userID uuid.UUID, ttl time.Duration) (*models.RefreshToken, error) {
	// Generate a random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	tokenString := hex.EncodeToString(tokenBytes)

	//Create UUID for the tokenID
	tokenID := uuid.New()
	expiresAt := time.Now().Add(ttl)

	token := &models.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     tokenString, // secure random. token
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetRefreshToken retrieves a refresh token by its token string
func (r *RefreshTokenRepository) GetRefreshToken(tokenString string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token = ?
	`

	var token models.RefreshToken
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

// // RevokeAllUserTokens revokes all refresh tokens for a specific user
// func (r *RefreshTokenRepository) RevokeAllUserTokens(userID uuid.UUID) error {
// 	query := `
// 		UPDATE refresh_tokens
// 		SET revoked = true
// 		WHERE user_id = ? AND revoked = false
// 	`

// 	_, err := r.db.Exec(query, userID)
// 	return err
// }

// // DeleteExpiredTokens removes all expired tokens from the database
// func (r *RefreshTokenRepository) DeleteExpiredTokens() error {
// 	query := `
// 		DELETE FROM refresh_tokens
// 		WHERE expires_at < ?
// 	`

// 	_, err := r.db.Exec(query, time.Now())
// 	return err
// }
