package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"
	"time"
)

func setupTestUserDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
        balance INTEGER NOT NULL DEFAULT 0
    )`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	form := "username=testuser&first_name=test&last_name=user&email=test@example.co.jp&password=secret&password_confirm=secret&date_of_birth=1993-09-14"

	//Build request
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	t.Run("POST request should return form data", func(t *testing.T) {

		//Capture response
		recorder := httptest.NewRecorder()

		//Call handler
		handleRegistrationForm(recorder, request)

		var resp User
		if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to parse response %v", err)
		}

		timeString := "1993-09-14"
		time, err := time.Parse(time.DateOnly, timeString)

		if err != nil {
			t.Errorf("error parsing time string %v", err)
		}

		want := User{Username: "testuser", FirstName: "test", LastName: "user", Email: "test@example.co.jp", Password: "secret", DateOfBirth: time}

		if resp.Username != want.Username {
			t.Errorf("wrong username input")
		}
		if resp.FirstName != want.FirstName {
			t.Errorf("wrong firstName input")
		}
		if resp.LastName != want.LastName {
			t.Errorf("wrong lastName input")
		}
		if resp.Email != want.Email {
			t.Errorf("wrong Email input")
		}
		if resp.Password != want.Password {
			t.Errorf("wrong password input")
		}
		if resp.DateOfBirth != want.DateOfBirth {
			t.Errorf("wrong DOB input, %+v", resp.DateOfBirth)
		}

		if resp.Username != "testuser" {
			t.Errorf("expected username 'testuser' but got %q,", resp.Username)
		}

		if _, err := mail.ParseAddress(resp.Email); err != nil {
			t.Errorf("Expected valid email address but got %q: %v", resp.Email, err)
		}

		if resp.FirstName == " " {
			t.Errorf("Expected valid first name field %v", resp.FirstName)
		}
	})

	t.Run("should create a new user from http response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handleRegistrationForm(recorder, request)
		var testUser User

		if err := json.NewDecoder(recorder.Body).Decode(&testUser); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		db := setupTestUserDB(t)
		defer db.Close()
		repo := &SQLiteUserRepo{db: db}

		if err := repo.CreateUser(&testUser); err != nil {
			t.Fatalf("failed to create user: %v\n, %+v", err, testUser)
		}

		if testUser.ID == 0 {
			t.Errorf("expected non-zero ID after insert")
		}

		got, err := repo.GetUser(testUser.ID)
		if err != nil {
			t.Fatalf("failed to get user:%v", err)
		}

		if got.Username != testUser.Username || got.Email != testUser.Email || got.Password != testUser.Password {
			t.Errorf("got %+v\n, want %+v,", got, testUser)
		}
	})
}

func TestPasswordHashing(t *testing.T) {
	password := "secret123"
	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == password {
		t.Error("Hash should not be the same as password")
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("Password check failed")
	}
}
