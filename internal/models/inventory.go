package models

import (
	"time"

	"github.com/google/uuid"
)

// InventoryItem represents an item in a user's inventory
type InventoryItem struct {
	UserID     uuid.UUID `json:"user_id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int       `json:"quantity"`
	AcquiredAt time.Time `json:"acquired_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// InventoryItemDetail includes product details for display
type InventoryItemDetail struct {
	InventoryItem
	Product Product `json:"product"`
}
