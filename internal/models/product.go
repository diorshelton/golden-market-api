package models

import (
	"time"
)

// Product represents a product in the marketplace
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       Coins     `json:"price"`
	Stock       int       `json:"stock"`
	RestockRate int       `json:"restock_rate"`
	MaxStock    int       `json:"max_stock"`
	VendorID    int       `json:"vendor_id"`
	LastRestock time.Time `json:"last_restock"`
	CreatedAt   time.Time `json:"created_at"`
}

// Item represents a product in a user's inventory
type Item struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// Coins represents the in-game currency
type Coins int32
