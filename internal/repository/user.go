package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *pgx.Conn
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database
func (r *UserRepository) CreateUser(username, firstName, lastName, email, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}

	query := `
	INSERT INTO users (id, username, first_name, last_name, email, password_hash, balance, created_at, last_login)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	ctx := context.Background()

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Balance,
		user.CreatedAt,
		user.LastLogin,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, password_hash, balance, created_at, last_login
		FROM users
		WHERE email = $1
	`

	var user models.User
	var lastLogin sql.NullTime

	ctx := context.Background()

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
		&lastLogin,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time.UTC()
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, password_hash, balance, created_at, last_login
		FROM users
		WHERE username = $1
	`

	var user models.User
	var lastLogin sql.NullTime

	ctx := context.Background()

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
		&lastLogin,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time.UTC()
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, password_hash, balance, created_at, last_login
		FROM users
		WHERE id = $1
	`

	var user models.User
	var lastLogin sql.NullTime

	ctx := context.Background()

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
		&lastLogin,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time.UTC()
	}

	return &user, nil
}

// UpdateBalance updates a user's coin balance
func (r *UserRepository) UpdateBalance(userID uuid.UUID, newBalance models.Coins) error {
	query := `UPDATE users SET balance = $1 WHERE id = $2`

	ctx := context.Background()

	_, err := r.db.Exec(ctx, query, newBalance, userID)
	return err
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(userID uuid.UUID) error {
	query := `UPDATE users SET last_login = $1 WHERE id = $2`

	ctx := context.Background()

	_, err := r.db.Exec(ctx, query, time.Now().UTC(), userID)
	return err
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT * FROM users`

	ctx := context.Background()

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.PasswordHash,
			&u.Balance,
			&u.CreatedAt,
			&u.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, rows.Err()
}
