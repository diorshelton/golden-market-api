package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("could not connect:%w", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Connected to PostgresSQL")
	return conn, nil
}

// SetupTestUserDB creates a temporary in memory test user database  with both users and refresh_tokens tables
func SetupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	usersQuery := `
		CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		balance INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_login DATETIME
		)
	`
	_, err = db.Exec(usersQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	refreshTokensQuery := `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		token TEXT NOT NULL UNIQUE,
		expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		revoked BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
	`
	_, err = db.Exec(refreshTokensQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	return db
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
			username VARCHAR(50) NOT NULL UNIQUE,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			balance INTEGER NOT NULL DEFAULT 0,
			INVENTORY JSONB[],
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_login TIMESTAMP WITH TIME ZONE
		)
	`
	_, err = db.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	return db, nil
}

// SetupTestRefreshTokenDB creates an in-memory database with both tables
// (refresh tokens need users table to exist for foreign key)
func SetupTestRefreshTokenDB() *sql.DB {
	// Just call SetupTestDB since refresh_tokens requires users table
	return SetupTestDB()
}

// CleanupTestDB closes and cleans up the test database
func CleanupTestDB(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			fmt.Printf("Warning: Failed to close test database: %v\n", err)
		}
	}
}
