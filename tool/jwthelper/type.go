package jwthelper

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// SecretProvider is the interface for providing the HMAC secret.
type SecretProvider interface {
	GetHmacSecret() []byte
}

// TokenType is the type of authentication.
// It can be either AccessToken or RefreshToken.
type TokenType string

var (
	AccessToken  TokenType = "AccessToken"
	RefreshToken TokenType = "RefreshToken"
)

// Config is the configuration for the JWT helper.
type Config struct {
	Issuer               string
	Audiences            []string
	AccessTokenValidity  time.Duration
	RefreshTokenValidity time.Duration
}

// Jwt is the implementation for the JWT helper.
type Jwt struct {
	hmacSecret           []byte
	issuer               string
	audiences            []string
	accessTokenValidity  time.Duration
	refreshTokenValidity time.Duration
	SigningMethod        *jwt.SigningMethodHMAC
}

// CustomClaims wraps custom claim `token_type` with jwt.RegisteredClaims for the JWT token.
type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}
