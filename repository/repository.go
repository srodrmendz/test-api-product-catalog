package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	internalErr "github.com/srodrmendz/api-auth/errors"
	"github.com/srodrmendz/api-auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Create new user repository
func New(client *mongo.Client, database string, collection string) *UsersRepository {
	return &UsersRepository{
		collection: client.Database(database).Collection(collection),
	}
}

// Create a new user
func (r *UsersRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	user.ID = uuid.NewString()

	// Set created and updated date at current date
	now := time.Now()

	user.CreatedAt = now

	user.UpdatedAt = now

	// Generate bcrypt hashed password
	password, err := r.generatePassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("generating password %w", err)
	}

	user.Password = *password

	// Create user on repository
	if _, err := r.collection.InsertOne(ctx, user); err != nil {
		writeError, ok := err.(mongo.WriteException)
		if !ok {
			return nil, fmt.Errorf("creating user %s on repository %w", user.Email, err)
		}

		for _, wErr := range writeError.WriteErrors {
			if wErr.Code == mongoDBDuplicatedKeyErrorCode {
				return nil, internalErr.ErrUserAlreadyExist
			}
		}

		return nil, fmt.Errorf("creating user %s on repository %w", user.Email, err)
	}

	return &user, nil
}

// Authenticate user credentials
func (r *UsersRepository) Authenticate(ctx context.Context, email string, password string) (*model.UserResponse, error) {
	response := r.collection.FindOne(ctx, bson.M{"email": email})

	// Check if user exist on repository
	if response.Err() != nil {
		return nil, internalErr.ErrUserNotFound
	}

	var user model.User

	if err := response.Decode(&user); err != nil {
		return nil, fmt.Errorf("decoding user from repository %w", err)
	}

	// Check user password is valid
	if !r.checkPasswordIsValid(user.Password, password) {
		return nil, internalErr.ErrUserNotFound
	}

	return model.MapUserToResponse(user), nil
}

// Delete user
func (r *UsersRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return fmt.Errorf("deleting user from repository %w", err)
	}

	return nil
}

// Get user by id
func (r *UsersRepository) GetByID(ctx context.Context, id string) (*model.UserResponse, error) {
	response := r.collection.FindOne(ctx, bson.M{"_id": id})

	// Check if user exist on repository
	if response.Err() != nil {
		return nil, internalErr.ErrUserNotFound
	}

	var user model.User

	if err := response.Decode(&user); err != nil {
		return nil, fmt.Errorf("decoding user from repository %w", err)
	}

	return model.MapUserToResponse(user), nil
}

// Generate user hashed password using bcrypt
func (r *UsersRepository) generatePassword(pass string) (*string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	password := string(b)

	return &password, nil
}

// Check user bcrypt hashed password with password
func (r *UsersRepository) checkPasswordIsValid(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	// Password is valid if error is nil
	return err == nil
}
