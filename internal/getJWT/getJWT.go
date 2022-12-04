package getJWT

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	Email string
	jwt.StandardClaims
}

const key = "some strings to identify my key"

func getJWT(email string) (string, error) {

	userClaims := claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		Email: email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &userClaims)

	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("couldn't Sign String: %w", err)
	}

	return signedString, nil
}
