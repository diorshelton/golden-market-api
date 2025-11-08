package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
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

	if err == sql.ErrNoRows {
		// Insert new item
		insertQuery := `
		INSERT INTO cart_items(id, user_id, product_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		`
		now := time.Now().UTC()
		_, err = r.db.Exec(insertQuery, uuid.New(), userID, productID, quantity, now, now)
		if err != nil {
			return fmt.Errorf("failed to add to cart: %w", err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to check existing cart item: %w", err)
	}

	// Updating existing item
	updateQuery := `
		UPDATE cart_items
		SET quantity = quantity + ?, updated_at = ?
		WHERE id = ?
	`
	_, err = r.db.Exec(updateQuery, quantity, time.Now().UTC(), existingID)
	if err != nil {
		return fmt.Errorf("failed to update cart quantity: %w", err)
	}

	return nil
}

// GetCart retrieves all items in a user's cart and product details
func (r *CartRepository) GetCart(userID uuid.UUID) (*models.CartSummary, error) {
	query := `
		SELECT
			ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
			p.id, p.name, p.description, p.price, p.stock, p.image_url, p.category, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = ?
		ORDER BY ci.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	defer rows.Close()
	var items []models.CartItem
	totalItems := 0
	totalPrice := 0

	for rows.Next() {
		var item models.CartItem
		var product models.Product
		var imageURL sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&imageURL,
			&product.Category,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		if imageURL.Valid {
			product.ImageURL = imageURL.String
		}

		item.Product = &product
		items = append(items, item)

		totalItems += item.Quantity
		totalPrice += int(product.Price) * item.Quantity
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart items: %w", err)
	}

	return &models.CartSummary{
		Items:      items,
		TotalItems: totalPrice,
		TotalPrice: totalPrice,
	}, nil
}

// UpdateCartItemQuantity updated the quantity of a specific cart item
func (r *CartRepository) UpdateCartItemQuantity(cartItem uuid.UUID, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, quantity, time.Now(), cartItem)
	if err != nil {
		return fmt.Errorf("failed to update cart item %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

func (r *CartRepository) RemoveFromCart(userID, cartItemID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, userID, cartItemID)
	if err != nil {
		return fmt.Errorf("failed to remove from cart %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart item not found %w", err)
	}

	return nil
}
