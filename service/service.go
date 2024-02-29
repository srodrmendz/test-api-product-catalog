package service

import (
	"context"

	"github.com/srodrmendz/api-product-catalog/model"
	"github.com/srodrmendz/api-product-catalog/repository"
)

// Create new product service
func New(repository repository.Repository) *ProductsCatalogService {
	return &ProductsCatalogService{
		repository: repository,
	}
}

// Create a new product
func (s *ProductsCatalogService) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	return s.repository.Create(ctx, product)
}

// Get a product by id
func (s *ProductsCatalogService) GetByID(ctx context.Context, id string) (*model.Product, error) {
	return s.repository.GetByID(ctx, id)
}

// Get a product by sku
func (s *ProductsCatalogService) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	return s.repository.GetBySKU(ctx, sku)
}

// Delete a product
func (s *ProductsCatalogService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

// Update a product
func (s *ProductsCatalogService) Update(ctx context.Context, request *model.Update) (*model.Product, error) {
	return s.repository.Update(ctx, request.ID, request.Qty)
}

// Search products
func (s *ProductsCatalogService) Search(ctx context.Context, request model.SearchRequest) (*model.SearchResponse, error) {
	products, total, err := s.repository.Search(ctx, request.Limit, request.Offset)
	if err != nil {
		return nil, err
	}

	return &model.SearchResponse{
		Products: products,
		Metadata: model.Metadata{
			Total:  *total,
			Limit:  request.Limit,
			Offset: request.Offset,
		},
	}, nil
}
