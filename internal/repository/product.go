package repository

import (
	"database/sql"
	"time"

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
		INSERT INTO products (name, description, price, stock, restock_rate, max_stock, vendor_id, last_restock, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.RestockRate,
		product.MaxStock,
		product.VendorID,
		product.LastRestock,
		time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int(id)
	return nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, restock_rate, max_stock, vendor_id, last_restock, created_at
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
		&product.RestockRate,
		&product.MaxStock,
		&product.VendorID,
		&product.LastRestock,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetByVendorID retrieves all products for a specific vendor
func (r *ProductRepository) GetByVendorID(vendorID int) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, restock_rate, max_stock, vendor_id, last_restock, created_at
		FROM products
		WHERE vendor_id = ?
	`

	rows, err := r.db.Query(query, vendorID)
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
			&product.RestockRate,
			&product.MaxStock,
			&product.VendorID,
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

// GetAll retrieves all products
func (r *ProductRepository) GetAll() ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, restock_rate, max_stock, vendor_id, last_restock, created_at
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
			&product.RestockRate,
			&product.MaxStock,
			&product.VendorID,
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
func (r *ProductRepository) UpdateStock(productID, newStock int) error {
	query := `UPDATE products SET stock = ? WHERE id = ?`
	_, err := r.db.Exec(query, newStock, productID)
	return err
}

// UpdateLastRestock updates the last restock timestamp
func (r *ProductRepository) UpdateLastRestock(productID int, lastRestock time.Time) error {
	query := `UPDATE products SET last_restock = ? WHERE id = ?`
	_, err := r.db.Exec(query, lastRestock, productID)
	return err
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(productID int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.Exec(query, productID)
	return err
}
