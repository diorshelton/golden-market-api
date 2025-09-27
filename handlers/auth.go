package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/diorshelton/golden-market/auth"
)

func HandleRegisterForm (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "An error occurred: %v", err)
	}
	// Get form values and trim white space
	username := strings.TrimSpace(r.Form.Get("username"))
	firstName := strings.TrimSpace(r.Form.Get("first_name"))
	lastName := strings.TrimSpace(r.Form.Get("last_name"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))
	passwordConfirm := strings.TrimSpace(r.Form.Get("password_confirm"))
	dobStr := strings.TrimSpace(r.Form.Get("date_of_birth"))

	// Check input for empty strings
	if username == "" || firstName == "" || lastName == "" || email == "" || password == "" || dobStr == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
	}
	// Verify passwords match
	if password != passwordConfirm {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
	}
	// Parse date_of_birth (format from <input type="date"> is YYYY-MM-DD)
	dob, err := time.Parse("2006-01-02", dobStr)

	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
	}
	// Validate email address
	if _, err := mail.ParseAddress(email); err != nil {
		http.Error(w, "invalid email address", http.StatusBadRequest)
	}

	password, err = auth.HashPassword(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	type input struct {
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
}

	user := input{
		Username:    username,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Password:    password,
		DateOfBirth: dob,
		CreatedAt:   time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}