package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *pgx.Conn
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) (models.Product, error) {
	query := `
		INSERT INTO products (id, name, description, price, stock, last_restock, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)

	`
	result, err := r.db.Exec(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	if err != nil {
		return *product, fmt.Errorf("failed to insert product:%v", err)
	}

	if result.RowsAffected() != 1 {
		return *product, fmt.Errorf("there were no rows affected: %v", err)
	}

	return *product, nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, restock_rate, max_stock, vendor_id, last_restock, created_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.LastRestock,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetAll retrieves all products
func (r *ProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, restock_rate, max_stock, vendor_id, last_restock, created_at
		FROM products
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.LastRestock,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, rows.Err()
}

// UpdateStock updates the stock level for a product
func (r *ProductRepository) UpdateStock(ctx context.Context, productID, newStock int) error {
	query := `UPDATE products SET stock = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, newStock, productID)
	return err
}

// UpdateLastRestock updates the last restock timestamp
func (r *ProductRepository) UpdateLastRestock(ctx context.Context, productID int, lastRestock time.Time) error {
	query := `UPDATE products SET last_restock = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, lastRestock, productID)
	return err
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(ctx context.Context, productID int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, productID)
	return err
}
