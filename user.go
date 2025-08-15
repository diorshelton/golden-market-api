package main

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	ID          int
	Username    string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Email       string
	Password    string
	Balance     Coins
	Inventory   []Item
	CreatedAt   time.Time
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}