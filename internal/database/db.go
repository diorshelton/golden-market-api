package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// SetupTestUserDB creates a temporary in memory test user database
func SetupTestUserDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Printf("Database Error:%v", err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS users (
		id string PRIMARY KEY,
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
	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Failed to create table: %v", err)
	}
	return db
}

// SetupTestUserDB creates a temporary in memory test refresh token database
func SetupRefreshTokenDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Printf("Database Error:%v", err)
	}

	query := `
		CREATE TABLE refresh_tokens (
		id string PRIMARY KEY,
		user_id TEXT NOT NULL UNIQUE,
		token TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		revoked BOOLEAN DEFAULT 0
		)
	`
	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Failed to create table: %v", err)
	}
	return db
}
