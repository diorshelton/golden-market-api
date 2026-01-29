package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// AddToCart adds a product to the user's cart or updates quantity it if exists
func (r *CartRepository) AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	// Check if item already exists in cart
	var existingID uuid.UUID
	var existingQty int

	query := `SELECT id, quantity FROM cart_items WHERE user_id = $1 AND product_id = $2`
	err := r.db.QueryRow(ctx, query, userID, productID).Scan(&existingID, &existingQty)

	if err == pgx.ErrNoRows {
		// Insert new item
		insertQuery := `
			INSERT INTO cart_items(id, user_id, product_id, quantity, added_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		now := time.Now().UTC()
		_, err = r.db.Exec(ctx, insertQuery, uuid.New(), userID, productID, quantity, now, now)
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
		SET quantity = quantity + $1, updated_at = $2
		WHERE id = $3
	`
	_, err = r.db.Exec(ctx, updateQuery, quantity, time.Now().UTC(), existingID)
	if err != nil {
		return fmt.Errorf("failed to update cart quantity: %w", err)
	}

	return nil
}

// GetCart retrieves all items in a user's cart and product details
func (r *CartRepository) GetCart(ctx context.Context, userID uuid.UUID) (*models.CartSummary, error) {
	query := `
		SELECT
			ci.id, ci.user_id, ci.product_id, ci.quantity, ci.added_at, ci.updated_at,
			p.id, p.name, p.description, p.price, p.stock, p.image_url, p.category, p.is_available, p.last_restock, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
		ORDER BY ci.added_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	defer rows.Close()

	var items []models.CartItemDetail
	totalItems := 0
	totalPrice := models.Coins(0)

	for rows.Next() {
		var cartItem models.CartItem
		var product models.Product
		var imageURL *string

		err := rows.Scan(
			&cartItem.ID,
			&cartItem.UserID,
			&cartItem.ProductID,
			&cartItem.Quantity,
			&cartItem.AddedAt,
			&cartItem.UpdatedAt,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&imageURL,
			&product.Category,
			&product.IsAvailable,
			&product.LastRestock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		if imageURL != nil {
			product.ImageURL = *imageURL
		}

		//Calculate subtotal for cartItem
		subtotal := models.Coins(int(product.Price) * cartItem.Quantity)

		//Build CartItemDetail with product info and subtotal
		itemDetail := models.CartItemDetail{
			CartItemID: cartItem.ID,
			Product:    product,
			Quantity:   cartItem.Quantity,
			Subtotal:   subtotal,
		}
		items = append(items, itemDetail)

		totalItems += cartItem.Quantity
		totalPrice += models.Coins(int(product.Price) * cartItem.Quantity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart items: %w", err)
	}

	return &models.CartSummary{
		Items:      items,
		TotalItems: totalItems,
		TotalPrice: models.Coins(totalPrice),
	}, nil
}

// UpdateCartItemQuantity updated the quantity of a specific cart item
func (r *CartRepository) UpdateCartItemQuantity(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := r.db.Exec(ctx, query, quantity, time.Now().UTC(), cartItemID)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// RemoveFromCart removes an item from the user's cart
func (r *CartRepository) RemoveFromCart(ctx context.Context, userID, cartItemID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, cartItemID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove from cart: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// ClearCart removes all items from a user's cart (within a transaction)
func (r *CartRepository) ClearCart(ctx context.Context, tx DBTX, userID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}
