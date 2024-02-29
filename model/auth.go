package model

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// Represents auth response when user is authenticated successfully
type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresIn int64     `json:"expires_in"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Represents user auth request
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Represents jwt claims of signed token
type JWTClaim struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}
