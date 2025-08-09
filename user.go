package main

import "time"

type User struct {
	ID        int
	Username string
	FirstName string
	LastName string
	DateOfBirth time.Time
	Email string
	Password string
	Balance Coins
	Inventory []Item
	CreatedAt time.Time
}
