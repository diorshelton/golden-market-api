package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(product *models.Product) error {

	query := `
	INSERT INTO products (id, name, description, price, stock, image_url, category, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStock updates the stock quantity for a product
func (r *ProductRepository) UpdateStock(productID uuid.UUID, newStock int) error {
	query := `UPDATE products SET stock = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, newStock, productID)
	if err != nil {
		return fmt.Errorf("failed to update stock %w", err)
	}
	return nil
}

func (r ProductRepository) DecrementStock(productID uuid.UUID, quantity int) error {
	query := `
	UPDATE products
	SET stock = stock - ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ? AND stock >= ?
	`
	result, err := r.db.Exec(query, quantity, productID, quantity)
	if err != nil {
		return fmt.Errorf("failed to decrement stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient stock or product not found")
	}
	return nil
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(productID int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.Exec(query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// GetAll retrieves all products
func (r *ProductRepository) GetAll(category string, minPrice, maxPrice int) ([]models.Product, error) {

	query := `
		SELECT id, name, description, price, stock, image_url,
		category, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	args := []any{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if minPrice > 0 {
		query += " AND price >= ?"
		args = append(args, minPrice)
	}

	if maxPrice > 0 {
		query += " AND price <= ?"
		args = append(args, maxPrice)
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var imageURL sql.NullString

		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		if imageURL.Valid {
			product.ImageURL = imageURL.String
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(id uuid.UUID) (*models.Product, error) {

	query := `
		SELECT id, name, description, price, stock, image_url, category,created_at, updated_at
		FROM products
		WHERE id = ?
	`

	var product models.Product
	var imageURL sql.NullString

	err := r.db.QueryRow(query, id).Scan(
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

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if imageURL.Valid {
		product.ImageURL = imageURL.String
	}

	return &product, nil
}

// SearchProducts searches product by name or description

func (r ProductRepository) SearchProducts(searchTerm string) ([]models.Product, error) {
	query := `
	SELECT id, name, description, price, stock, image_url, category, created_at, updated_at
	FROM products
	WHERE name LIKE ? OR description LIKE ?
	ORDER BY name ASC
 `

	searchPattern := "%" + strings.ToLower(searchTerm) + "%"
	rows, err := r.db.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search products %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var imageURL sql.NullString

		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		if imageURL.Valid {
			product.ImageURL = imageURL.String
		}

		products = append(products, product)
	}
	return products, nil
}
