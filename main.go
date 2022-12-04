package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	http.HandleFunc("/", serveHtml)
	http.HandleFunc("/submit", submitInfo)
	http.ListenAndServe(":8080", nil)
}

func getJWT(message string) (string, error) {
	key := "some strings"

	type claims struct {
		jwt.StandardClaims
		Email string
	}

	userClaims := claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		Email: message,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &userClaims)

	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("couldn't Sign String: %w", err)
	}

	return signedString, nil
}

func submitInfo(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	email := request.FormValue("email")
	if email == "" {
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	signedString, err := getJWT(email)
	if err != nil {
		http.Error(response, "couldn't get JWT", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:  "session",
		Value: signedString,
	}

	http.SetCookie(response, &cookie)
	http.Redirect(response, request, "/", http.StatusSeeOther)
}

func serveHtml(response http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("session")
	if err != nil {
		cookie = &http.Cookie{}
	}

	isEqual := true

	message := "Not logged in"
	if isEqual {
		message = "Logged in"
	}

	html := `
	<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>Doc Cookie</title>
		</head>
		<body>
      <p>Cookie value: ` + cookie.Value + `</p>
      <p>` + message + `</p>
			<form action="/submit" method="post">
				<input type="email" name="email" />
				<input type="submit" />
			</form>
		</body>
	</html>`

	io.WriteString(response, html)
}
