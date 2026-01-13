package models

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the marketplace
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       Coins     `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url,omitempty"`
	Category    string    `json:"category"`
	IsAvailable bool      `json:"is_available"`
	LastRestock time.Time `json:"last_restock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Item represents a product in a user's inventory
type Item struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// Coins represents in-game currency
type Coins int64
