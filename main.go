package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

func main() {
	s := &http.Server{
		Addr:    ":3000",
		Handler: marketServer(),
	}

	log.Printf("Running sever on port%v", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func marketServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /sign-in", handleSignIn)
	mux.HandleFunc("GET /register", handleRegister)
	mux.HandleFunc("POST /register", handleRegistrationForm)
	return mux
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, ".public/sign-in.html")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/register.html")
}

func handleRegistrationForm(w http.ResponseWriter, r *http.Request) {
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

	if _, err := mail.ParseAddress(email); err != nil {
		http.Error(w, "invalid email address", http.StatusBadRequest)
	}

	password, err = HashPassword(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user := User{
		Username:    username,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Password:    password,
		DateOfBirth: dob,
		CreatedAt:   time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
