package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        balance INTEGER DEFAULT 0,
				email TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL UNIQUE
    );`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestSQLiteUserRepo_CreateAndGetUser(t *testing.T) {
	db := setupTestDB(t)
	repo := &SQLiteUserRepo{db: db}

	user := &User{Username: "tomoyo"}
	if err := repo.CreateUser(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Errorf("expected non-zero ID after insert")
	}

	got, err := repo.GetUser(user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if got.Username != user.Username || got.Balance != user.Balance {
		t.Errorf("got %+v, want %+v", got, user)
	}
}
