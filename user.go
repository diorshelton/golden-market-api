package main

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Balance     Coins     `json:"balance"`
	Inventory   []Item    `json:"inventory,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
