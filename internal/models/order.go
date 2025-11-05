package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the state of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a completed purchase
type Order struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	TotalAmount int         `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Item        []OrderItem `json:"items"`
}

// OrderItem represents a single product in an order
type OrderItem struct {
	ID           uuid.UUID `json:"id"`
	OrderID      uuid.UUID `json:"order_id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	Quantity     int       `json:"quantity"`
	PricePerUnit int       `json:"price_per_unit"`
	Subtotal     int       `json:"subtotal"`
}
