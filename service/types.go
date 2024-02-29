package service

import (
	"context"
	"time"

	"github.com/srodrmendz/api-auth/model"
	"github.com/srodrmendz/api-auth/repository"
)

// Check on build time that AuthService implement Service interface
var _ Service = (*AuthService)(nil)

const expiresAt = 1 * time.Hour

// Service defines the methods that should be implemented by a auth service.
type Service interface {
	// Authenticate a user
	// Returns error if user is not found or there is an error in the system
	Authenticate(ctx context.Context, email string, password string) (*model.AuthResponse, error)
}

// Service Implementation
type AuthService struct {
	repository   repository.Repository
	jwtSecretKey string
}
