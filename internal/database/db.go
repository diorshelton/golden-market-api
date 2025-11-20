package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// SetupTestUserDB creates a temporary in memory test user database  with both users and refresh_tokens tables
func loadEnv(envString string) string {
// Try to load .env file, but don't crash if it doesn't exist
 if	err := godotenv.Load(".env"); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	dbString := os.Getenv(envString)
	if dbString == "" {
		log.Fatalf("%s not set in env",envString)
	}
	return dbString
}

func SetupTestDB() (*pgx.Conn, error) {
	dbString := loadEnv("TEST_DB_URL")

	ctx := context.Background()

	db, err := pgx.Connect(ctx, dbString)
	if err != nil {
		log.Fatalf("Failed to open test database: %v", err)
	}
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

	refreshTokensQuery := `
	CREATE TEMPORARY TABLE refresh_tokens (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(512) NOT NULL UNIQUE,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		revoked BOOLEAN NOT NULL DEFAULT FALSE
		);
	`
	_, err = db.Exec(ctx, refreshTokensQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh_tokens table: %v", err)
	}

	productsQuery := `
	CREATE TEMPORARY TABLE products(
	  id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (name <>''),
    description TEXT,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    category VARCHAR(255),
    last_restock TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`
	_, err = db.Exec(ctx, productsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create products table %v", err)
	}

	inventoryQuery := `
	CREATE TEMPORARY TABLE inventory (
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		quantity INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (user_id, product_id),
		CONSTRAINT quantity_positive CHECK (quantity > 0)
	);

		CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
		CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
		CREATE INDEX idx_inventory_user_id ON inventory(user_id);
		CREATE INDEX idx_inventory_product_id ON inventory(product_id);
		CREATE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_username ON users(username);
	`
	_, err = db.Exec(ctx, inventoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory table %v", err)
	}

	return db, nil
}

// SetupTestUserDB creates a temporary in-memory database with just the users table
// Use this only to test users in isolation
func SetupTestUserDB() (*pgx.Conn, error) {
	connString := os.Getenv("TEST_DATABASE_URL")
	if connString == " " {
		return nil, fmt.Errorf("TEST_DATABASE_URL not set")
	}

	ctx := context.Background()

	db, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to open test database: %v", err)
	}

	query := `
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
		)
	`
	_, err = db.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	return db, nil
}

// CleanupTestDB closes and cleans up the test database
func CleanupTestDB(db *pgx.Conn) {
	ctx := context.Background()
	if db != nil {
		if err := db.Close(ctx); err != nil {
			fmt.Printf("Warning: Failed to close test database: %v\n", err)
		}
	}
}
