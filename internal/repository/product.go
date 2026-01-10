package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *pgxpool.Pool
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	now := time.Now().UTC()

	product.CreatedAt = now
	product.UpdatedAt = now
	product.LastRestock = now

	query := `
		INSERT INTO products (id, name, description, price, stock, image_url, category, last_restock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
		product.Category,
		product.LastRestock,
		product.IsAvailable,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}

// UpdateStock updates the stock quantity for a product
func (r *ProductRepository) UpdateStock(ctx context.Context, productID uuid.UUID, newStock int) error {
	query := `UPDATE products SET stock = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, newStock, productID)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}
	return nil
}

// DecrementStock decrements the stock by the given quantity
func (r *ProductRepository) DecrementStock(ctx context.Context, productID uuid.UUID, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1, updated_at = NOW()
		WHERE id = $2 AND stock >= $3
	`
	result, err := r.db.Exec(ctx, query, quantity, productID, quantity)
	if err != nil {
		return fmt.Errorf("failed to decrement stock: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock or product not found")
	}
	return nil
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(ctx context.Context, productID uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// GetAll retrieves all products with optional filtering
func (r *ProductRepository) GetAll(ctx context.Context, category string, minPrice, maxPrice int) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, image_url, category, last_restock, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	args := []any{}
	paramCount := 1

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", paramCount)
		args = append(args, category)
		paramCount++
	}

	if minPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", paramCount)
		args = append(args, minPrice)
		paramCount++
	}

	if maxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", paramCount)
		args = append(args, maxPrice)
		paramCount++
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		var imageURL *string

		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		if imageURL != nil {
			product.ImageURL = *imageURL
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, image_url, category, last_restock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	var imageURL *string

	err := r.db.QueryRow(ctx, query, id).Scan(
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
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if imageURL != nil {
		product.ImageURL = *imageURL
	}

	return &product, nil
}
