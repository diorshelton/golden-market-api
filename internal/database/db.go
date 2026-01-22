package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// SetupTestUserDB creates a temporary in memory test user database  with both users and refresh_tokens tables
func loadEnv(envString string) string {
	// Try to load .env file from multiple possible locations
	// This handles both running from root and from test subdirectories
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	dbString := os.Getenv(envString)
	if dbString == "" {
		log.Fatalf("%s not set in env", envString)
	}
	return dbString
}

func SetupTestDB() (*pgxpool.Pool, error) {
	dbString := loadEnv("TEMP_DB_URL")

	ctx := context.Background()

	db, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Fatalf("Failed to open test database: %v", err)
	}

	// Create users table
	usersQuery := `
	CREATE TEMPORARY TABLE users (
		id UUID PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		balance INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		last_login TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, usersQuery)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create refresh_tokens table
	refreshTokensQuery := `
	CREATE TEMPORARY TABLE refresh_tokens (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(512) NOT NULL UNIQUE,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		revoked BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err = db.Exec(ctx, refreshTokensQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh_tokens table: %v", err)
	}

	// Create products table
	productsQuery := `
	CREATE TEMPORARY TABLE products (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL CHECK (name <> ''),
		description TEXT,
		price INTEGER NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0,
		image_url TEXT,
		category VARCHAR(255),
		is_available BOOLEAN NOT NULL DEFAULT true,
		last_restock TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, productsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create products table: %v", err)
	}

	// Create inventory table
	inventoryQuery := `
	CREATE TEMPORARY TABLE inventory (
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
		quantity INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (user_id, product_id),
		CONSTRAINT quantity_positive CHECK (quantity > 0)
	);`

	_, err = db.Exec(ctx, inventoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory table: %v", err)
	}

	// Create cart_items table
	cartItemsQuery := `
	CREATE TEMPORARY TABLE cart_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		UNIQUE(user_id, product_id)
	);`

	_, err = db.Exec(ctx, cartItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart_items table: %v", err)
	}

	// Create all indexes
	indexQuery := `
		CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
		CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
		CREATE INDEX idx_inventory_user_id ON inventory(user_id);
		CREATE INDEX idx_inventory_product_id ON inventory(product_id);
		CREATE INDEX idx_cart_items_user_id ON cart_items(user_id);
		CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);
		CREATE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_username ON users(username);
	`

	_, err = db.Exec(ctx, indexQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %v", err)
	}

	return db, nil
}

func SetupDB() (*pgxpool.Pool, error) {
	dbString := loadEnv("LOCAL_DB_URL")
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dbString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create users table
	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		balance INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		last_login TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, usersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	// Create refresh_tokens table
	refreshTokensQuery := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(512) NOT NULL UNIQUE,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		revoked BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err = db.Exec(ctx, refreshTokensQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh_tokens table: %w", err)
	}

	// Create products table
	productsQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL CHECK (name <> ''),
		description TEXT,
		price INTEGER NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0,
		image_url TEXT,
		category VARCHAR(255),
		is_available BOOLEAN NOT NULL DEFAULT true,
		last_restock TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, productsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create products table: %w", err)
	}

	// Create inventory table
	inventoryQuery := `
	CREATE TABLE IF NOT EXISTS inventory (
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
		quantity INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (user_id, product_id),
		CONSTRAINT quantity_positive CHECK (quantity > 0)
	);`

	_, err = db.Exec(ctx, inventoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory table: %w", err)
	}

	// Create cart_items table
	cartItemsQuery := `
	CREATE TABLE IF NOT EXISTS cart_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		UNIQUE(user_id, product_id)
	);`

	_, err = db.Exec(ctx, cartItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart_items table: %w", err)
	}

	// Create all indexes
	indexQuery := `
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
		CREATE INDEX IF NOT EXISTS idx_inventory_user_id ON inventory(user_id);
		CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory(product_id);
		CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
		CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`

	_, err = db.Exec(ctx, indexQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return db, nil
}
