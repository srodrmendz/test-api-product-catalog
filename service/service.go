package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/srodrmendz/api-auth/model"
	"github.com/srodrmendz/api-auth/repository"
)

// Create new auth service
func New(repository repository.Repository, jwtSecretKey string) *AuthService {
	return &AuthService{
		repository:   repository,
		jwtSecretKey: jwtSecretKey,
	}
}

// Authenticate user
func (s *AuthService) Authenticate(ctx context.Context, email string, password string) (*model.AuthResponse, error) {
	// First authenticate email and password on repository
	user, err := s.repository.Authenticate(ctx, email, password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	expires := now.Add(expiresAt)

	// Create JWT
	token, err := s.generateJwt(user.Email, user.UserName, now, expires)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token:     *token,
		ExpiresIn: expires.Unix(),
		ExpiresAt: expires,
	}, nil
}

// Generate JWT
func (s *AuthService) generateJwt(email string, username string, issuedAt, expires time.Time) (*string, error) {
	// First create claims
	claims := model.JWTClaim{
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expires.Unix(),
			Id:        uuid.NewString(),
			Issuer:    "test_app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tk, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return nil, fmt.Errorf("signing jwt token %w", err)
	}

	return &tk, nil
}
