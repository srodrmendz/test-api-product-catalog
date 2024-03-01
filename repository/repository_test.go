package repository

import (
	"context"
	"fmt"
	"testing"

	internalError "github.com/srodrmendz/api-product-catalog/errors"
	"github.com/srodrmendz/api-product-catalog/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test Repository
func TestRepository_Integration_Create(t *testing.T) {
	t.Run("successfully create product", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		product, err := repo.Create(ctx, model.Product{
			Qty:    10,
			Name:   "product 1",
			Sku:    "sku1",
			Images: []string{"https://google.com/image1.png"},
		})

		// Then
		require.NoError(t, err)

		response, err := repo.GetByID(ctx, product.ID)

		require.NoError(t, err)

		assert.Equal(t, product.ID, response.ID)

		assert.Equal(t, product.Qty, response.Qty)

		assert.Equal(t, true, product.InStock)

		assert.Equal(t, false, product.CreatedAt.IsZero())

		assert.Equal(t, false, product.UpdatedAt.IsZero())

		err = repo.Delete(ctx, response.ID)

		require.NoError(t, err)
	})

	t.Run("failed to create product, sku already founded", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		product, err := repo.Create(ctx, model.Product{
			Qty:    10,
			Name:   "product 1",
			Sku:    "sku1",
			Images: []string{"https://google.com/image1.png"},
		})

		// Then
		require.NoError(t, err)

		_, err = repo.Create(ctx, model.Product{
			Qty:  100,
			Name: "product 2",
			Sku:  product.Sku,
		})

		assert.EqualError(t, internalError.ErrProductSKUAlreadyExist, err.Error())

		err = repo.Delete(ctx, product.ID)

		require.NoError(t, err)
	})
}

func TestRepository_Integration_GetByID(t *testing.T) {
	t.Run("successfully get product by id", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		product, err := repo.Create(ctx, model.Product{
			Qty:    10,
			Name:   "product 1",
			Sku:    "sku1",
			Images: []string{"https://google.com/image1.png"},
		})

		require.NoError(t, err)

		response, err := repo.GetByID(ctx, product.ID)

		// Then
		require.NoError(t, err)

		assert.Equal(t, product.ID, response.ID)

		assert.Equal(t, product.Qty, response.Qty)

		assert.Equal(t, true, product.InStock)

		assert.Equal(t, false, product.CreatedAt.IsZero())

		assert.Equal(t, false, product.UpdatedAt.IsZero())

		err = repo.Delete(ctx, response.ID)

		require.NoError(t, err)
	})

	t.Run("failed to get by id, product not found", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		response, err := repo.GetByID(ctx, "fake")

		// Then
		assert.Error(t, internalError.ErrProductNotFound, err.Error())

		assert.Nil(t, response)
	})
}

func TestRepository_Integration_GetBySKU(t *testing.T) {
	t.Run("successfully get product by sku", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		product, err := repo.Create(ctx, model.Product{
			Qty:    10,
			Name:   "product 1",
			Sku:    "sku1",
			Images: []string{"https://google.com/image1.png"},
		})

		require.NoError(t, err)

		response, err := repo.GetBySKU(ctx, product.Sku)

		// Then
		require.NoError(t, err)

		assert.Equal(t, product.ID, response.ID)

		assert.Equal(t, product.Qty, response.Qty)

		assert.Equal(t, true, product.InStock)

		assert.Equal(t, false, product.CreatedAt.IsZero())

		assert.Equal(t, false, product.UpdatedAt.IsZero())

		err = repo.Delete(ctx, response.ID)

		require.NoError(t, err)
	})

	t.Run("failed to get by sku, product not found", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		response, err := repo.GetBySKU(ctx, "fake")

		// Then
		assert.Error(t, internalError.ErrProductNotFound, err.Error())

		assert.Nil(t, response)
	})
}

func TestRepository_Integration_Delete(t *testing.T) {
	t.Run("successfully delete product", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		product, err := repo.Create(ctx, model.Product{
			Qty:    10,
			Name:   "product 1",
			Sku:    "sku1",
			Images: []string{"https://google.com/image1.png"},
		})

		// Then
		require.NoError(t, err)

		err = repo.Delete(ctx, product.ID)

		require.NoError(t, err)
	})
}

func TestRepository_Integration_Update(t *testing.T) {
	t.Run("successfully update product", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		product, err := repo.Create(ctx, model.Product{
			Name:   "product 1",
			Sku:    "sku1",
			Images: []string{"https://google.com/image1.png"},
		})

		require.NoError(t, err)

		resp, err := repo.Update(ctx, product.ID, 1000)

		// Then
		assert.NoError(t, err)

		assert.NotNil(t, resp)

		assert.Equal(t, true, resp.InStock)

		assert.Equal(t, uint64(1000), resp.Qty)

		err = repo.Delete(ctx, product.ID)

		require.NoError(t, err)
	})

	t.Run("failed update product, product not found", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		product, err := repo.Update(ctx, "fake", 1000)

		// Then
		require.Error(t, err)

		assert.Error(t, internalError.ErrProductNotFound, err.Error())

		assert.Nil(t, product)
	})
}

func TestRepository_Integration_Search(t *testing.T) {
	t.Run("successfully search products", func(t *testing.T) {
		// Given
		ctx := context.TODO()

		client := createMongoClient(t, ctx)

		defer client.Disconnect(ctx)

		repo := New(client, "ecommerce", "products")

		// When
		productsIDS := make([]string, 0)

		for i := 0; i < 5; i++ {
			product, err := repo.Create(ctx, model.Product{
				Qty:    10,
				Name:   "product 1",
				Sku:    fmt.Sprintf("sku%d", i),
				Images: []string{"https://google.com/image1.png"},
			})

			require.NoError(t, err)

			productsIDS = append(productsIDS, product.ID)
		}

		resp, total, err := repo.Search(ctx, "", "asc", true, 0, 0)

		// Then
		require.NoError(t, err)

		require.NoError(t, err)

		assert.Less(t, 0, len(resp))

		assert.Less(t, 0, int(*total))

		for _, id := range productsIDS {
			err = repo.Delete(ctx, id)

			require.NoError(t, err)
		}
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
