package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
	Password string
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "sign-in")
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
	email := r.Form.Get("e-mail")
	password := r.Form.Get("password")

	user := NewUser{
		Username: username,
		Email:    email,
		Message:  "registration successful",
		Password: password,
	}

	json.NewEncoder(w).Encode(user)
	// fmt.Println(w, userData)
}

func marketServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /sign-in", handleSignIn)
	mux.HandleFunc("GET /register", handleRegister)
	mux.HandleFunc("POST /register", handleRegistrationForm)
	return mux
}

func main() {
	s := &http.Server{
		Addr:    ":3000",
		Handler: marketServer(),
	}
	log.Printf("Running sever on port%v", s.Addr)
	log.Fatal(s.ListenAndServe())
}
