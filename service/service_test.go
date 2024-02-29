package service

import (
	"context"
	"testing"
	"time"

	internalErr "github.com/srodrmendz/api-auth/errors"
	"github.com/srodrmendz/api-auth/model"
	"github.com/srodrmendz/api-auth/repository"
	"github.com/srodrmendz/api-auth/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Generate JWT Method
func TestService_GenerateJWT(t *testing.T) {
	// Given
	secretKey := "test-key"

	srv := New(nil, secretKey)

	email := "fake@email.com"

	username := "fake"

	// When
	token, err := srv.generateJwt(email, username, time.Now(), time.Now().Add(expiresAt))

	// Then
	require.NoError(t, err)

	assert.NotNil(t, token)

	claims, err := utils.GetClaimsFromToken(*token, secretKey)

	require.NoError(t, err)

	assert.NotNil(t, claims)

	assert.Equal(t, username, claims.Username)

	assert.Equal(t, email, claims.Email)
}

// Test Authenticate method
func TestService_Authenticate(t *testing.T) {
	dataTable := []struct {
		name          string
		repository    repository.Repository
		expectedError error
	}{
		{
			name: "failed to create authentication response, user credentials are invalid",
			repository: &mockRepository{
				err: internalErr.ErrUserNotFound,
			},
			expectedError: internalErr.ErrUserNotFound,
		},
		{
			name: "successfully create authentication response",
			repository: &mockRepository{
				user: &model.UserResponse{},
			},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			srv := New(dt.repository, "")

			// When
			resp, err := srv.Authenticate(context.TODO(), "", "")

			// Then
			if err != nil {
				assert.Equal(t, dt.expectedError, err)

				return
			}

			require.NotNil(t, resp)

			assert.Equal(t, false, resp.ExpiresAt.IsZero())

			assert.NotEqual(t, 0, resp.ExpiresIn)

			assert.NotEqual(t, "", resp.Token)
		})
	}
}

type mockRepository struct {
	err  error
	user *model.UserResponse
}

func (m *mockRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	return nil, nil
}

func (m *mockRepository) Authenticate(ctx context.Context, email string, password string) (*model.UserResponse, error) {
	return m.user, m.err
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*model.UserResponse, error) {
	return nil, nil
}
