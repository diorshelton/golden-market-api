package main

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
				date_of_birth DATE NOT NULL,
        username TEXT NOT NULL,
				first_name TEXT NOT NULL,
				last_name TEXT NOT NULL,
        balance INTEGER DEFAULT 0,
				email TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
				created_at DATE NOT NULL
    );`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestSQLiteUserRepo_CreateAndGetUser(t *testing.T) {
	db := setupTestDB(t)
	repo := &SQLiteUserRepo{db: db}

	//Test User
	var user = &User{Username: "tomoyo1", FirstName: "Tomoyo", LastName: "Daidouji", Email: "daito@gmail.com.jp"}

	// Convert time from string to time.Time type
	dob, err := time.Parse("2006-01-02", "1993-09-03")
	if err != nil {
		t.Errorf("An error occurred: %v", err)
	}

	//Add to struct
	user.DateOfBirth = dob
	user.CreatedAt = time.Now().UTC()

	user.Password, err = HashPassword("password")
	if err != nil {
		t.Errorf("An error occurred while hashing password: %v", err)
	}

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

	if got.Username == "" || got.FirstName == "" || got.LastName == "" || got.Email == "" || got.Password == "" {
		t.Fatalf("User missing value %+v", got)
	}

	if got.DateOfBirth.IsZero() {
		t.Errorf("Date of birth not set %v, want %v", got.DateOfBirth, dob)
	}

	if got.CreatedAt.IsZero() {
		t.Errorf("Date of CreatedAt field not set, got: %v", got.CreatedAt)
	}
}
