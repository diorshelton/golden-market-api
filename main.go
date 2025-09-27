package main

import (
	"log"
	"net/http"

	"github.com/diorshelton/golden-market/handlers"
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
	mux.HandleFunc("POST /register", handlers.HandleRegisterForm)
	return mux
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, ".public/sign-in.html")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/register.html")
}
