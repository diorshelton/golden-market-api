package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	username := r.Form.Get("username")
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	dobStr := r.Form.Get("date_of_birth")
	dob, err := time.Parse("2006-01-02", dobStr)

	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
	}

	user := User{
		Username:    username,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Password:    password,
		DateOfBirth: dob,
	}

	json.NewEncoder(w).Encode(user)
}
