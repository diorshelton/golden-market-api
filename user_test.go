package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	// "net/mail"
	// "reflect"
	"strings"
	"testing"
	// "time"
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
	registrationInput := "username=testuser&first_name=test&last_name=user&email=test@example.co.jp&password=secret&password_confirm=secret&date_of_birth=1993-09-14"

	//Declare request
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(registrationInput))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// t.Run("POST request should return form data", func(t *testing.T) {

	// 	//Capture response
	// 	recorder := httptest.NewRecorder()

	// 	//Call handler
	// 	handleRegistrationForm(recorder, request)

	//	//Decode json response
	// 	var resp User
	// 	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
	// 		t.Fatalf("failed to parse response %v", err)
	// 	}

	// 	dobString := "1993-09-14"
	// 	dob, err := time.Parse(time.DateOnly, dobString)

	// 	if err != nil {
	// 		t.Fatalf("error parsing time string %v", err)
	// 	}

	// 	want := User{Username: "testuser", FirstName: "test", LastName: "user", Email: "test@example.co.jp", Password: "secret", DateOfBirth: dob}

	// 	if !reflect.DeepEqual(want, resp) {
	// 		t.Fatalf("\nStructs don't match\n got: %v\nwant: %v", want, resp)
	// 	}

	// 	if _, err := mail.ParseAddress(resp.Email); err != nil {
	// 		t.Errorf("Expected valid email address but got %q: %v", resp.Email, err)
	// 	}
	// })

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
	t.Run("username and email must be unique", func(t *testing.T) {

		formInput := "username=adventure1&first_name=finn&last_name=mertens&email=finMertens@advntrtm.com&password=enchiridion&password_confirm=enchiridion&date_of_birth=1993-09-14"

		request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(formInput))

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		recorder := httptest.NewRecorder()

		handleRegistrationForm(recorder, request)

		var finn User

		if err := json.NewDecoder(recorder.Body).Decode(&finn); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		db := setupTestUserDB(t)
		defer db.Close()
		repo := &SQLiteUserRepo{db: db}

		if err := repo.CreateUser(&finn); err != nil {
			t.Fatalf("failed to create user: %v\n, %+v", err, finn)
		}

		got, err := repo.GetUser(finn.ID)
		if err != nil {
			t.Fatalf("failed to get user:%v", err)
		}
		t.Errorf("Got user: %v", got)
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
