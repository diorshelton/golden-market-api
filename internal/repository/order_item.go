package repository

import (
	"context"
	"fmt"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItemRepository struct {
	db *pgxpool.Pool
}

func NewOrderItemRepository(db *pgxpool.Pool) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

// Create inserts a new order item into the database (within a transaction)
func (r *OrderItemRepository) Create(ctx context.Context, tx DBTX, item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, product_name, quantity, price_per_unit, subtotal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.Exec(ctx, query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.ProductName,
		item.Quantity,
		item.PricePerUnit,
		item.Subtotal,
		item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}

	return nil
}

// GetByOrderID retrieves all items for an order
func (r *OrderItemRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, product_name, quantity, price_per_unit, subtotal, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.PricePerUnit,
			&item.Subtotal,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order items: %w", err)
	}

	return items, nil
}
