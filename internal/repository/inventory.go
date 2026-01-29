package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// AddOrUpdate performs an UPSERT to add items to inventory or update quantity if exists
func (r *InventoryRepository) AddOrUpdate(ctx context.Context, tx DBTX, userID, productID uuid.UUID, quantity int) error {
	now := time.Now().UTC()

	query := `
		INSERT INTO inventory (user_id, product_id, quantity, acquired_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = inventory.quantity + $3, updated_at = $5
	`

	_, err := tx.Exec(ctx, query, userID, productID, quantity, now, now)
	if err != nil {
		return fmt.Errorf("failed to add/update inventory: %w", err)
	}

	return nil
}

// GetByUserID retrieves all inventory items for a user with product details
func (r *InventoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.InventoryItemDetail, error) {
	query := `
		SELECT
			i.user_id, i.product_id, i.quantity, i.acquired_at, i.updated_at,
			p.id, p.name, p.description, p.price, p.stock, p.image_url, p.category, p.is_available, p.last_restock, p.created_at, p.updated_at
		FROM inventory i
		JOIN products p ON i.product_id = p.id
		WHERE i.user_id = $1
		ORDER BY i.acquired_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItemDetail
	for rows.Next() {
		var item models.InventoryItemDetail
		var imageURL *string

		err := rows.Scan(
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
			&item.AcquiredAt,
			&item.UpdatedAt,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Price,
			&item.Product.Stock,
			&imageURL,
			&item.Product.Category,
			&item.Product.IsAvailable,
			&item.Product.LastRestock,
			&item.Product.CreatedAt,
			&item.Product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}

		if imageURL != nil {
			item.Product.ImageURL = *imageURL
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating inventory: %w", err)
	}

	return items, nil
}

// ClearByUserID removes all inventory items for a user (admin function)
func (r *InventoryRepository) ClearByUserID(ctx context.Context, tx DBTX, userID uuid.UUID) error {
	query := `DELETE FROM inventory WHERE user_id = $1`

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to clear inventory: %w", err)
	}

	return nil
}
