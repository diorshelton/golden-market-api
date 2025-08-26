package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "sign-in")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r, "./public/register.html")
}

func handleRegistrationForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	email := r.Form.Get("e-mail")
	password := r.Form.Get("password")
	fmt.Fprintf(w,"Email: %v, PW:%v, Username:%v,",email, password, username)
}

func marketServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /sign-in", handleSignIn)
	mux.HandleFunc("GET /register", handleRegister)
	mux.HandleFunc("POST /register", handleRegistrationForm)
	return  mux
}

func main() {
	s := &http.Server{
		Addr: ":3000",
		Handler: marketServer(),
	}
	log.Printf("Running sever on port%v", s.Addr)
	log.Fatal(s.ListenAndServe())
}
