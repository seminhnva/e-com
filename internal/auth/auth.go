package auth

import "github.com/golang-jwt/jwt/v5"

// auth
type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
