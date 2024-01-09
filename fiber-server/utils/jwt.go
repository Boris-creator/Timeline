package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateTokenString(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return t
}

func ParseTokenString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return token, err
	}
	if !token.Valid {
		return token, errors.New("token is invalid")
	}
	return token, nil
}
