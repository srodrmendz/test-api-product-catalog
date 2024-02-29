package service

import (
	"context"
	"errors"
	"testing"

	internalErrors "github.com/srodrmendz/api-product-catalog/errors"
	"github.com/srodrmendz/api-product-catalog/model"
	"github.com/srodrmendz/api-product-catalog/repository"
	"github.com/stretchr/testify/assert"
)

// Test Create method
func TestService_Create(t *testing.T) {
	dataTable := []struct {
		name        string
		repository  repository.Repository
		expectedErr error
	}{
		{
			name:        "failed to create product, error on repository",
			expectedErr: internalErrors.ErrProductSKUAlreadyExist,
			repository: &mockRepository{
				err: internalErrors.ErrProductSKUAlreadyExist,
			},
		},
		{
			name: "successfully create product",
			repository: &mockRepository{
				product: &model.Product{},
			},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			srv := New(dt.repository)

			// When
			product, err := srv.Create(context.TODO(), model.Product{
				Name: "Test",
			})

			// Then
			if err != nil {
				assert.EqualError(t, dt.expectedErr, err.Error())

				return
			}

			assert.NotEmpty(t, product.Name)
		})
	}
}

// Test Get By ID method
func TestService_GetByID(t *testing.T) {
	dataTable := []struct {
		name        string
		repository  repository.Repository
		expectedErr error
	}{
		{
			name:        "failed to get product by id, product not found",
			expectedErr: internalErrors.ErrProductNotFound,
			repository: &mockRepository{
				err: internalErrors.ErrProductNotFound,
			},
		},
		{
			name: "successfully get product by id",
			repository: &mockRepository{
				product: &model.Product{
					Name: "Product",
				},
			},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			srv := New(dt.repository)

			// When
			product, err := srv.GetByID(context.TODO(), "")

			// Then
			if err != nil {
				assert.EqualError(t, dt.expectedErr, err.Error())

				return
			}

			assert.NotEmpty(t, product.Name)
		})
	}
}

// Test Get By SKU method
func TestService_GetBySKU(t *testing.T) {
	dataTable := []struct {
		name        string
		repository  repository.Repository
		expectedErr error
	}{
		{
			name:        "failed to get product by sku, product not found",
			expectedErr: internalErrors.ErrProductNotFound,
			repository: &mockRepository{
				err: internalErrors.ErrProductNotFound,
			},
		},
		{
			name: "successfully get product by sku",
			repository: &mockRepository{
				product: &model.Product{
					Name: "Product",
				},
			},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			srv := New(dt.repository)

			// When
			product, err := srv.GetBySKU(context.TODO(), "")

			// Then
			if err != nil {
				assert.EqualError(t, dt.expectedErr, err.Error())

				return
			}

			assert.NotEmpty(t, product.Name)
		})
	}
}

// Test Delete method
func TestService_Delete(t *testing.T) {
	dataTable := []struct {
		name        string
		repository  repository.Repository
		expectedErr error
	}{
		{
			name:        "failed to delete product",
			expectedErr: errors.New("error on repository"),
			repository: &mockRepository{
				err: errors.New("error on repository"),
			},
		},
		{
			name:       "successfully delete product",
			repository: &mockRepository{},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			srv := New(dt.repository)

			// When
			err := srv.Delete(context.TODO(), "")

			// Then
			assert.Equal(t, dt.expectedErr, err)
		})
	}
}

// Test Update method
func TestService_Update(t *testing.T) {
	dataTable := []struct {
		name        string
		repository  repository.Repository
		expectedErr error
	}{
		{
			name:        "failed to update product, error on repository",
			expectedErr: internalErrors.ErrProductNotFound,
			repository: &mockRepository{
				err: internalErrors.ErrProductNotFound,
			},
		},
		{
			name: "successfully update product",
			repository: &mockRepository{
				product: &model.Product{
					Name: "Product 1",
				},
			},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			srv := New(dt.repository)

			// When
			product, err := srv.Update(context.TODO(), &model.Update{})

			// Then
			if err != nil {
				assert.EqualError(t, dt.expectedErr, err.Error())

				return
			}

			assert.NotEmpty(t, product.Name)
		})
	}
}

// Test Search method
func TestService_Search(t *testing.T) {
	dataTable := []struct {
		name          string
		repository    repository.Repository
		expectedErr   error
		expectedTotal int64
	}{
		{
			name:        "failed to search products, error on repository",
			expectedErr: errors.New("error on repository"),
			repository: &mockRepository{
				err: errors.New("error on repository"),
			},
		},
		{
			name: "successfully search products",
			repository: &mockRepository{
				total: 100,
			},
			expectedTotal: 100,
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			srv := New(dt.repository)

			request := model.SearchRequest{
				Limit:  10,
				Offset: 0,
			}

			// When
			response, err := srv.Search(context.TODO(), request)

			// Then
			if err != nil {
				assert.EqualError(t, dt.expectedErr, err.Error())

				return
			}

			assert.Equal(t, dt.expectedTotal, response.Metadata.Total)

			assert.Equal(t, request.Limit, response.Metadata.Limit)

			assert.Equal(t, request.Offset, response.Metadata.Offset)
		})
	}
}

type mockRepository struct {
	err      error
	product  *model.Product
	products []model.Product
	total    int64
}

func (m *mockRepository) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	return &product, m.err
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	return m.product, m.err
}

func (m *mockRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	return m.product, m.err
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	return m.err
}

func (m *mockRepository) Update(ctx context.Context, id string, qty uint64) (*model.Product, error) {
	return m.product, m.err
}

func (m *mockRepository) Search(ctx context.Context, limit int, offset int) ([]model.Product, *int64, error) {
	return m.products, &m.total, m.err
}
