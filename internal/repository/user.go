package repository

import (
	"database/sql"
	"time"

	"github.com/diorshelton/golden-market/internal/models"
	"github.com/google/uuid"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(
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
		WHERE email = ?
	`

	var user models.User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
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
		WHERE id = ?
	`

	var user models.User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
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

// // UpdateBalance updates a user's coin balance
// func (r *UserRepository) UpdateBalance(userID uuid.UUID, newBalance models.Coins) error {
// 	query := `UPDATE users SET balance = ? WHERE id = ?`
// 	_, err := r.db.Exec(query, newBalance, userID)
// 	return err
// }

// // UpdateLastLogin updates the user's last login timestamp
// func (r *UserRepository) UpdateLastLogin(userID uuid.UUID) error {
// 	query := `UPDATE users SET last_login = ? WHERE id = ?`
// 	_, err := r.db.Exec(query, time.Now().UTC(), userID)
// 	return err
// }
