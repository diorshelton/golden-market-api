package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user := User{
		Username: username,
		Email:    email,
		Password: password,
	}

	json.NewEncoder(w).Encode(user)
}
