package repository

import (
	"database/sql"

	"github.com/diorshelton/golden-market-api/internal/models"
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

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {

	query := `
		SELECT id, name, description, price, stock, image_url,  category, created_at, updated_at
		FROM products
		WHERE id = ?
	`

	var product models.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.Category,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetAll retrieves all products
func (r *ProductRepository) GetAll() ([]*models.Product, error) {

	query := `
		SELECT id, name, description, price, stock, image_url,
		category, created_at, updated_at
		FROM products
	`

	rows, err := r.db.Query(query)
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
			&product.ImageURL,
			&product.Category,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, rows.Err()
}

// UpdateStock updates the stock level for a product
func (r *ProductRepository) UpdateStock(productID, newStock int) error {
	query := `UPDATE products SET stock = ? WHERE id = ?`
	_, err := r.db.Exec(query, newStock, productID)
	return err
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(productID int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.Exec(query, productID)
	return err
}
