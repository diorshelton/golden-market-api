package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/diorshelton/golden-market/auth"
	"github.com/diorshelton/golden-market/db"
	"github.com/diorshelton/golden-market/models"
)

func main() {
	s := &http.Server{
		Addr:    ":3000",
		Handler: marketServer(),
	}
	log.Printf("Running sever on port%v", s.Addr)
	log.Fatal(s.ListenAndServe())
}

var userDatabase = db.SetupTestUserDB()
var refreshTokenDb = db.SetupRefreshTokenDB()

var userRepo = models.NewUserRepository(userDatabase)
var tokenRepo = models.NewRefreshTokenRepository(refreshTokenDb)

// var service = auth.NewAuthService()

func marketServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /sign-in", handleSignIn)
	mux.HandleFunc("GET /register", handleRegister)
	// mux.HandleFunc("POST /register", handlers.HandleRegisterForm)
	return mux
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, ".public/sign-in.html")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/register.html")
}
