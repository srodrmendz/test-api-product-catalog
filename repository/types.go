package repository

import (
	"context"

	"github.com/srodrmendz/api-auth/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB error code when email index entry is duplicated
const mongoDBDuplicatedKeyErrorCode = 11000

// Check on build time that UsersRepository implement Repository interface
var _ Repository = (*UsersRepository)(nil)

// Repository defines the methods that should be implemented by a user repository.
type Repository interface {
	// Create a new user
	// Returns error if email already exists or there is an error in the system
	Create(ctx context.Context, user model.User) (*model.User, error)

	// Authenticate user with credentials
	// Returns error if credentials are invalid or there is an error in the system
	Authenticate(ctx context.Context, email string, password string) (*model.UserResponse, error)

	// Delete user
	// Returns error if there is an error in the system
	Delete(ctx context.Context, id string) error

	// Get user by id
	// Returns error if user does not exist or there is an error in the system
	GetByID(ctx context.Context, id string) (*model.UserResponse, error)
}

// MongoDB Users Repository Implementation
type UsersRepository struct {
	collection *mongo.Collection
}
