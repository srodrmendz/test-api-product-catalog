package repository

import (
	"context"
	"testing"

	internalErr "github.com/srodrmendz/api-auth/errors"
	"github.com/srodrmendz/api-auth/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test Repository
func TestRepository_Integration(t *testing.T) {
	t.Run("successfully create user", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "auth")

		password := "testPass"

		// When
		user, err := repo.Create(ctx, model.User{
			Email:    "fake@gmail.com",
			UserName: "fake",
			Password: password,
		})

		// Then
		require.NoError(t, err)

		response, err := repo.GetByID(ctx, user.ID)

		require.NoError(t, err)

		assert.Equal(t, user.ID, response.ID)

		assert.Equal(t, user.Email, response.Email)

		assert.Equal(t, user.NickName, response.NickName)

		assert.Equal(t, false, user.CreatedAt.IsZero())

		assert.Equal(t, false, user.UpdatedAt.IsZero())

		assert.Equal(t, true, repo.checkPasswordIsValid(user.Password, password))

		err = repo.Delete(ctx, response.ID)

		require.NoError(t, err)
	})

	t.Run("failed to create user, email already exist", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "auth")

		email := "fake@gmail.com"

		// When
		user, err := repo.Create(ctx, model.User{
			Email:    email,
			UserName: "fake",
			Password: "testPass",
		})

		// Then
		require.NoError(t, err)

		_, err = repo.Create(ctx, model.User{
			Email:    email,
			UserName: "fake",
			Password: "testPass",
		})

		assert.EqualError(t, internalErr.ErrUserAlreadyExist, err.Error())

		err = repo.Delete(ctx, user.ID)

		require.NoError(t, err)
	})
}

func createMongoClient(t *testing.T, ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("connecting mongo %s", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("connecting mongo %s", err)
	}

	return client
}
