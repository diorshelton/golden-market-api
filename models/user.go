package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Balance      Coins     `json:"balance"`
	Inventory    []Item    `json:"inventory,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login"`
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database
func (r *UserRepository) CreateUser(username, firstName, lastName, email, passwordHash string) (*User, error) {

	user := &User{
		ID:           uuid.New(),
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}

	query := `
		INSERT INTO users (id, username, first_name, last_name, email, password_hash, balance, created_at, last_login) VALUES (?,?,?,?,?,?,?,?,?)
		`
	_, err := r.db.Exec(query, user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Balance, user.CreatedAt, user.LastLogin)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, first_name, last_name, email, password_hash, balance, created_at, last_login FROM users WHERE email = ?`

	var user User
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
		&user.LastLogin,
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
func (r *UserRepository) GetUserByID(id uuid.UUID) (*User, error) {
	query := `SELECT id, username, first_name, last_name, email, password_hash, balance, created_at, last_login FROM users WHERE id = ?`

	var user User
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
		&user.LastLogin,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time.UTC()
	}

	return &user, nil
}
