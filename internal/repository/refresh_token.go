package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// RefreshTokenRepository handles database operations for refresh tokens
type RefreshTokenRepository struct {
	db *pgx.Conn
}

func NewRefreshTokenRepository(db *pgx.Conn) *RefreshTokenRepository {
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
		Token:     tokenString,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	ctx := context.Background()

	_, err := r.db.Exec(
		ctx, query,
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
		WHERE token = $1
	`

	var token models.RefreshToken

	ctx := context.Background()

	err := r.db.QueryRow(ctx, query, tokenString).Scan(
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

func (r *RefreshTokenRepository) DeleteRefreshToken(tokenString string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`

	ctx := context.Background()

	_, err := r.db.Exec(ctx, query, tokenString)
	return err
}

// RevokeRefreshToken marks a refresh token as revoked
func (r *RefreshTokenRepository) RevokeRefreshToken(tokenString string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE token = $1
	`

	ctx := context.Background()

	_, err := r.db.Exec(ctx, query, tokenString)
	return err
}

// DeleteExpiredTokens removes all expired tokens from the database
func (r *RefreshTokenRepository) DeleteExpiredTokens() error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < $1
	`

	ctx := context.Background()

	_, err := r.db.Exec(ctx, query, time.Now())
	return err
}

// RevokeAllUserTokens revokes all refresh tokens for a specific user
func (r *RefreshTokenRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE user_id = $1 AND revoked = false
	`

	ctx := context.Background()

	_, err := r.db.Exec(ctx, query, userID)
	return err
}
