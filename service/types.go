package service

import (
	"context"

	"github.com/srodrmendz/api-product-catalog/model"
	"github.com/srodrmendz/api-product-catalog/repository"
)

// Check on build time that Product Catalog implement Service interface
var _ Service = (*ProductsCatalogService)(nil)

// Service defines the methods that should be implemented by a Product Catalog Service.
type Service interface {
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
	Update(ctx context.Context, request *model.Update) (*model.Product, error)

	// Search products
	// Returns error if there is an error in the system
	Search(ctx context.Context, request model.SearchRequest) (*model.SearchResponse, error)
}

// Service Implementation
type ProductsCatalogService struct {
	repository repository.Repository
}
