package repository

import (
	"context"
	"testing"

	"github.com/diorshelton/golden-market-api/internal/database"
	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
)

type ProductData struct {
	Name        string
	Description string
	Price       int
	Stock       int
}

func TestProductRepository(t *testing.T) {

	dbConnection, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("An error occurred: %v", err)
	}

	products := []ProductData{
		{Name: "Mechanical Keyboard", Description: "RGB backlit keyboard with cherry MX brown switches, durable and responsive for typing and gaming.", Price: 12000, Stock: 40},
		{Name: "", Description: "27-inch 4K UHD monitor with USB-C connectivity and thin bezels, perfect for professional work.", Price: 39999, Stock: 0},
		{Name: "Noise-Cancelling Headphones", Description: "Over-ear headphones with active noise cancellation, 30-hour battery life, and superior sound quality.", Price: 19999, Stock: 200},
		{Name: "USB 3.0 Hub", Description: "A 4-port USB 3.0 hub with a dedicated power adapter for connecting multiple peripherals simultaneously.", Price: 2550, Stock: 300},
		{Name: "Wireless Ergonomic Mouse", Description: "A comfortable mouse with a 2.4 GHz wireless connection and an ergonomic design for long-term use.", Price: 4599, Stock: 150},
	}

	productRepo := NewProductRepository(dbConnection)

	t.Run("Create Product Happy Path", func(t *testing.T) {

		ctx := context.Background()

		keyboard := productCreator(products[0])

		err := productRepo.Create(ctx, keyboard)

		if err != nil {
			t.Errorf("Got an error %v", err)
		}

		if keyboard.Name == "" {
			t.Error("Empty string for name")
		}

	})

	t.Run("Create product without name", func(t *testing.T) {

		ctx := context.Background()
		noName := productCreator(products[1])

		err := productRepo.Create(ctx, noName)

		if err == nil {
			t.Error("Should have gotten and error but didn't", err)
		}

		results, err := queryProductsTable(ctx, *productRepo)

		t.Logf("Products found %v\n", results)
		if err != nil {
			t.Errorf("An error occurred whiled grabbing table: %v", err)
		}

	})
}

func productCreator(d ProductData) *models.Product {
	var product models.Product
	product.ID = uuid.New()
	product.Name = d.Name
	product.Description = d.Description
	product.Price = models.Coins(d.Price)
	product.Stock = d.Price

	return &product
}

func queryProductsTable(ctx context.Context, repo ProductRepository) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, image_url, category, last_restock, created_at, updated_at
		FROM products
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productList []models.Product
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
			&product.LastRestock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if imageURL != nil {
			product.ImageURL = *imageURL
		}
		productList = append(productList, product)
	}

	return productList, rows.Err()
}
