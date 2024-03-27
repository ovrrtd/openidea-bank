package common

import "github.com/golang-jwt/jwt/v5"

type ctxKey string

const (
	JwtCtxKey            ctxKey = "jwtContextKey"
	EncodedUserJwtCtxKey ctxKey = "encodedUserJwtCtxKey"
)

func (c ctxKey) ToString() string {
	return string(c)
}

type UserClaims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

type Meta struct {
	Limit  int
	Offset int
	Total  int
}

// regex

const (
	RegexEmailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)
