package jwt

import (
	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(payload jwt.Claims, secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJwt(tokenString string, claims jwt.Claims, secret string) error {
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if !tkn.Valid {
		return errorer.ErrUnauthorized
	}

	return nil
}
