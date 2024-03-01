package repository

import (
	"context"

	"github.com/srodrmendz/api-product-catalog/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB error code when sku index entry is duplicated
const mongoDBDuplicatedKeyErrorCode = 11000

// Check on build time that ProductsCatalogRepository implement Repository interface
var _ Repository = (*ProductsCatalogRepository)(nil)

// Repository defines the methods that should be implemented by a user repository.
type Repository interface {
	// Create a new product
	// Returns error if sku already exists or there is an error in the system
	Create(ctx context.Context, product model.Product) (*model.Product, error)

	// Get a product by id
	// Returns error if there is an error in the system
	GetByID(ctx context.Context, id string) (*model.Product, error)

	// Get a product by sku
	// Returns error if there is an error in the system
	GetBySKU(ctx context.Context, sku string) (*model.Product, error)

	// Delete a product
	// Returns error if there is an error in the system
	Delete(ctx context.Context, id string) error

	// Update a product
	// Returns error if there is an error in the system
	Update(ctx context.Context, id string, qty uint64) (*model.Product, error)

	// Search products
	// Returns error if there is an error in the system
	Search(ctx context.Context, name string, sort string, inStock bool, limit int, offset int) ([]model.Product, *int64, error)
}

// MongoDB Products Catalog Repository Implementation
type ProductsCatalogRepository struct {
	collection *mongo.Collection
}
