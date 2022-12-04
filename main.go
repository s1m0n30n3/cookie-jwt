package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	Email string
	jwt.StandardClaims
}

const key = "some strings to identify my key"

func main() {
	http.HandleFunc("/", servehtml)
	http.HandleFunc("/submit", submit)
	http.ListenAndServe(":8080", nil)
}
