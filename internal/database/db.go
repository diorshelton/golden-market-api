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
		balance INTEGER NOT NULL DEFAULT 5000,
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
		acquired_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		PRIMARY KEY (user_id, product_id),
		CONSTRAINT quantity_non_negative CHECK (quantity >= 0)
	);`

	_, err = db.Exec(ctx, inventoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory table: %v", err)
	}

	// Create orders table
	ordersQuery := `
	CREATE TEMPORARY TABLE orders (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		order_number VARCHAR(20) NOT NULL UNIQUE,
		total_amount INTEGER NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'completed',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, ordersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create orders table: %v", err)
	}

	// Create order_items table
	orderItemsQuery := `
	CREATE TEMPORARY TABLE order_items (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
		product_name VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL,
		price_per_unit INTEGER NOT NULL,
		subtotal INTEGER NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, orderItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create order_items table: %v", err)
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
		CREATE INDEX idx_orders_user_id ON orders(user_id);
		CREATE INDEX idx_orders_order_number ON orders(order_number);
		CREATE INDEX idx_orders_created_at ON orders(created_at);
		CREATE INDEX idx_order_items_order_id ON order_items(order_id);
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
		balance INTEGER NOT NULL DEFAULT 5000 CHECK (balance >= 0),
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
		stock INTEGER NOT NULL CHECK (stock >= 0),
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
		acquired_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		PRIMARY KEY (user_id, product_id),
		CONSTRAINT quantity_non_negative CHECK (quantity >= 0)
	);`

	_, err = db.Exec(ctx, inventoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory table: %w", err)
	}

	// Create orders table
	ordersQuery := `
	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		order_number VARCHAR(20) NOT NULL UNIQUE,
		total_amount INTEGER NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'completed',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, ordersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create orders table: %w", err)
	}

	// Create order_items table
	orderItemsQuery := `
	CREATE TABLE IF NOT EXISTS order_items (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
		product_name VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL,
		price_per_unit INTEGER NOT NULL,
		subtotal INTEGER NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err = db.Exec(ctx, orderItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create order_items table: %w", err)
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
		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
		CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
		CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
		CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
	`

	_, err = db.Exec(ctx, indexQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	// Run migrations for existing databases
	if err := runMigrations(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations handles schema updates for existing databases
func runMigrations(ctx context.Context, db *pgxpool.Pool) error {
	// Migration 1: Fix inventory constraint and add columns
	_, _ = db.Exec(ctx, `ALTER TABLE inventory DROP CONSTRAINT IF EXISTS quantity_positive`)
	_, _ = db.Exec(ctx, `ALTER TABLE inventory ADD CONSTRAINT quantity_non_negative CHECK (quantity >= 0)`)
	_, _ = db.Exec(ctx, `ALTER TABLE inventory ADD COLUMN IF NOT EXISTS acquired_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()`)
	_, _ = db.Exec(ctx, `ALTER TABLE inventory ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()`)

	// Migration 2: Update default balance for users (doesn't affect existing users)
	_, _ = db.Exec(ctx, `ALTER TABLE users ALTER COLUMN balance SET DEFAULT 5000`)

	return nil
}
