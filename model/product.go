package model

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Represent product structure
type Product struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description *string   `json:"description,omitempty" bson:"description,omitempty"`
	Sku         string    `json:"sku" bson:"sku"`
	Qty         uint64    `json:"qty" bson:"qty"`
	Images      []string  `json:"images,omitempty" bson:"images"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	Price       int64     `json:"price" bson:"price"`
	InStock     bool      `json:"in_stock" bson:"in_stock"`
}

// Represent search structure
type SearchResponse struct {
	Products []Product `json:"products"`
	Metadata Metadata  `json:"metadata"`
}

// Represent metadata structure used for pagination
type Metadata struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// Represent product update request
type UpdateRequest struct {
	Qty uint64 `json:"qty"`
}

// Represent product update
type Update struct {
	ID string
	UpdateRequest
}

// Represent search request
type SearchRequest struct {
	Limit   int
	Offset  int
	Name    string
	InStock bool
	Sort    string
}

// Build product create request and validate all requested data
type CreateBuilder struct {
	r *http.Request
}

func NewCreateBuilder(r *http.Request) *CreateBuilder {
	return &CreateBuilder{
		r: r,
	}
}

func (b *CreateBuilder) Build() (*Product, error) {
	var product Product

	if err := json.NewDecoder(b.r.Body).Decode(&product); err != nil {
		return nil, errors.New("incorrect product create body format")
	}

	if product.Name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if product.Sku == "" {
		return nil, errors.New("product sku cannot be empty")
	}

	for _, image := range product.Images {
		if image == "" {
			return nil, errors.New("product image cannot be empty")
		}
	}

	if product.Price <= 0 {
		return nil, errors.New("product price invalid value")
	}

	return &product, nil
}

// Build product update request and validate all requested data
type UpdateBuilder struct {
	r *http.Request
}

func NewUpdateBuilder(r *http.Request) *UpdateBuilder {
	return &UpdateBuilder{
		r: r,
	}
}

func (b *UpdateBuilder) Build() (*Update, error) {
	id := mux.Vars(b.r)["id"]
	if id == "" {
		return nil, errors.New("product id must be provided")
	}

	var request UpdateRequest

	if err := json.NewDecoder(b.r.Body).Decode(&request); err != nil {
		return nil, errors.New("incorrect product update body format")
	}

	return &Update{
		ID:            id,
		UpdateRequest: request,
	}, nil
}

// Build product search request and validate all requested data
type SearchBuilder struct {
	r *http.Request
}

func NewSearchBuilder(r *http.Request) *SearchBuilder {
	return &SearchBuilder{
		r: r,
	}
}

func (b *SearchBuilder) Build() (*SearchRequest, error) {
	limit, err := strconv.Atoi(b.r.URL.Query().Get("limit"))
	if err != nil {
		return nil, errors.New("incorrect limit format")
	}

	if limit < 0 {
		return nil, errors.New("incorrect limit format")
	}

	offset, err := strconv.Atoi(b.r.URL.Query().Get("offset"))
	if err != nil {
		return nil, errors.New("incorrect offset format")
	}

	if offset < 0 {
		return nil, errors.New("incorrect offset format")
	}

	inStock, _ := strconv.ParseBool(b.r.URL.Query().Get("in_stock"))

	return &SearchRequest{
		Limit:   limit,
		Offset:  offset,
		Name:    b.r.URL.Query().Get("name"),
		InStock: inStock,
		Sort:    b.r.URL.Query().Get("sort"),
	}, nil
}
