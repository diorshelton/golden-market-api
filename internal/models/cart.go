package models

import (
	"time"

	"github.com/google/uuid"
)

// CartItem represents a product in a user's shopping cart
type CartItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Product   *Product  `json:"product"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartSummary provides cart totals
type CartSummary struct {
	Items      []CartItem `json:"cart_item"`
	TotalItems int        `json:"total_items"`
	TotalPrice int        `json:"total_price"`
}
