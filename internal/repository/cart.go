package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

// AddToCart adds a product to the user's cart or updates quantity it if exists
func (r *CartRepository) AddToCart(userID, productID uuid.UUID, quantity int) error {
	// Check if item already exists in cart
	var existingID uuid.UUID
	var existingQty int

	query := `SELECT id, quantity FROM cart_items WHERE user_id = ? AND product_id = ?`
	err := r.db.QueryRow(query, userID, productID).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows{
		// Insert new item
		insertQuery := `
		INSERT INTO cart_items(id, user_id, product_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		`
		now := time.Now().UTC()
		_, err = r.db.Exec(insertQuery, uuid.New(), userID, productID, quantity, now, now)
		if err != nil  {
			return fmt.Errorf("failed to add to cart: %w", err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to check existing cart item: %w", err)
	}

	// Updating existing item
	updateQuery :=`
		UPDATE cart_items
		SET quantity = quantity + ?, updated_at = ?
		WHERE id = ?
	`
	_, err = r.db.Exec(updateQuery, quantity, time.Now().UTC(),existingID)
	if err != nil {
		return fmt.Errorf("failed to update cart quantity: %w", err)
	}

	return  nil
}
