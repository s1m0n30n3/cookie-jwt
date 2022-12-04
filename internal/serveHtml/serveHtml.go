package servehtml

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	Email string
	jwt.StandardClaims
}

const key = "some strings to identify my key"

func serveHtml(response http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("session")
	if err != nil {
		cookie = &http.Cookie{}
	}

	signedString := cookie.Value
	afterVerificationToken, err := jwt.ParseWithClaims(
		signedString,
		&claims{},
		func(beforeVerificationToken *jwt.Token) (interface{}, error) {
			if beforeVerificationToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("Someone tried to hack changed signing method")
			}
			return []byte(key), nil
		},
	)

	isEqual := err == nil && afterVerificationToken.Valid

	message := "Not logged in"
	if isEqual {
		message = "Logged in"
		userClaims := afterVerificationToken.Claims.(*claims)
		fmt.Println(userClaims.Email)
		fmt.Println(userClaims.StandardClaims.ExpiresAt)
		fmt.Println(userClaims.ExpiresAt)
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
