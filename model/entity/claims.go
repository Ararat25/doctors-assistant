package entity

import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
	Email string `json:"e"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	Email         string `json:"e"`
	AccessTokenID string `json:"aid"`
	jwt.RegisteredClaims
}
