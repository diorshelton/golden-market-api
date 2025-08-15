package main

import (
	"testing"
)

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
