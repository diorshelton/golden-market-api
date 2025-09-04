package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {

	form := "username=testuser&e-mail=test@example.co.jp&password=secret"

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

		if resp.Username != "testuser" {
			t.Errorf("expected username 'testuser' but got %q,", resp.Username)
		}

		email, err := mail.ParseAddress(resp.Email)

		if err != nil {
			log.Printf("Expected valid email address but got %v", email)
			log.Fatal(err)
		}

	})

	t.Run("should create a new user", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handleRegistrationForm(recorder, request)

		var resp User

		if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		CreateNewUser(&resp)

		if len(userStore) = 0 {
			t.Errorf("Expected a new user from %+v but did not get one", resp)
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
