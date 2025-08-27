package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("Retrieves form data", func(t *testing.T) {
		form := "username=testuser&e-mail=test@example.com&password=secret"

		//Build request
		request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		//Capture response
		recorder := httptest.NewRecorder()

		//Call handler
		handleRegistrationForm(recorder, request)

		var resp NewUser

		if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to parse response %v", err)
		}

		if resp.Message != "registration successful" {
			t.Errorf("expected username %q, got %q", "testuser", resp.Username)
		}
	})
	t.Run("Email not already being used", func(t *testing.T) {

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
