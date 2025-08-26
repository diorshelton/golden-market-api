package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	server := marketServer()
	t.Run("Create new user", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/register", nil)
		recorder := httptest.NewRecorder()

		server.ServeHTTP(recorder,request)

		got := recorder.Body.String()
		want := "sign-in"

		if got != want {
			t.Errorf("Got %v, want %v", got, want)
		}
	})
	t.Run("Email not already being used", func(t *testing.T) {

	// 	request, _ := http.NewRequest(http.MethodPost, "/register", nil)
	// 	recorder := httptest.NewRecorder()

	// 	server.ServeHTTP(recorder,request)

	// 	got := recorder.Body.String()
	// 	want := "sign-in"

	// 	if got != want {
	// 		t.Errorf("Got %v but want %v",got, want)
	// 	}
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
