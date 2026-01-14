package models

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	AddedAt   time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItemDetail struct {
	CartItemID uuid.UUID `json:"cart_item_id"`
	Product    Product   `json:"product"`
	Quantity   int       `json:"quantity"`
	Subtotal   Coins     `json:"subtotal"`
}

type CartSummary struct {
	Items      []CartItemDetail `json:"items"`
	TotalItems int              `json:"total_items"`
	TotalPrice Coins            `json:"total_price"`
}
