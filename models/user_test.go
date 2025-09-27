package models

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"


	_ "github.com/mattn/go-sqlite3"
	"github.com/diorshelton/golden-market/auth"
	"github.com/diorshelton/golden-market/handlers"
)

func TestCreateUser(t *testing.T) {
	registrationInput := "username=testuser&first_name=test&last_name=user&email=test@example.co.jp&password=secret&password_confirm=secret&date_of_birth=1993-09-14"

	//Declare request
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(registrationInput))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	db := setupTestUserDB(t)

	defer db.Close()
	repo := &SQLiteUserRepo{db: db}

	t.Run("POST request to form registration handler should return correct data", func(t *testing.T) {
		//Capture response
		recorder := httptest.NewRecorder()

		//Call handler
		handlers.HandleRegisterForm(recorder, request)

		//Decoded json response
		var resp User
		if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to parse response %v", err)
		}

		dobString := "1993-09-14"
		dob, err := time.Parse(time.DateOnly, dobString)

		if err != nil {
			t.Fatalf("error parsing time string %v", err)
		}

		//Declare comparable struct with testable fields
		want := User{Username: "testuser", FirstName: "test", LastName: "user", Email: "test@example.co.jp", Password: "secret", DateOfBirth: dob}

		if want.Username != resp.Username {
			t.Fatalf("\nStructs don't match\n got: %v\nwant: %v", resp, want)
		}

		if want.FirstName != resp.FirstName {
			t.Fatalf("FirstName not equal: want %v, got %v", want.FirstName, resp.FirstName)
		}

		if want.LastName != resp.LastName {
			t.Fatalf("Last Name not equal want %v got %v", want.LastName, resp.LastName)
		}

		if want.Email != resp.Email {
			t.Fatalf("Email not equal want %v got %v", want.Email, resp.Email)
		}

		if want.DateOfBirth != resp.DateOfBirth {
			t.Fatalf("DOB not equal want %v got %v", want.DateOfBirth, resp.DateOfBirth)
		}
	})

	t.Run("requests should create a new user in database", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handlers.HandleRegisterForm(recorder, request)

		var testUser User

		if err := json.NewDecoder(recorder.Body).Decode(&testUser); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		//Create User in database
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

		formInput2 :=
			"username=testuser&first_name=finn&last_name=mertens&email=finMertens@advntrtm.com&password=enchiridion&password_confirm=enchiridion&date_of_birth=1993-09-14"

		request2 := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(formInput2))

		request2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		recorder2 := httptest.NewRecorder()

		handlers.HandleRegisterForm(recorder2, request2)

		var testUser2 User

		if err := json.NewDecoder(recorder2.Body).Decode(&testUser2); err != nil {
			t.Fatalf("failed to parse response %v", err)
		}

		if err := repo.CreateUser(&testUser2); err == nil {
			t.Errorf("Did not receive duplicate username error but should have:%v", testUser2)
		}
	})
}

func TestPasswordHashing(t *testing.T) {
	password := "secret123"
	hash, err := auth.HashPassword(password)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == password {
		t.Error("Hash should not be the same as password")
	}

	 err = auth.VerifyPassword(password, hash); if err != nil {
		t.Error("Password check failed")
	}
}

func setupTestUserDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
				first_name TEXT NOT NULL,
				last_name TEXT NOT NULL,
				date_of_birth DATE NOT NULL,
				email TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
        balance INTEGER NOT NULL DEFAULT 0,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}
