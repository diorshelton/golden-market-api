package models

import (
	"database/sql"
	"time"

)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Balance     Coins     `json:"balance"`
	Inventory   []Item    `json:"inventory,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository (db * sql.DB) *UserRepository {
	return  &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, username, password string) (*User, error) {

	return &User{}, nil

}