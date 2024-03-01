package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	internalErrors "github.com/srodrmendz/api-product-catalog/errors"
	"github.com/srodrmendz/api-product-catalog/model"
	"github.com/srodrmendz/api-product-catalog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Healthcheck endpoint
func TestServer_HealthCheck(t *testing.T) {
	// Given
	version := "v1.0.0"

	buildDate := "02/29/2024"

	app := New(
		nil,
		mux.NewRouter(),
		"",
		version,
		buildDate,
	)

	endpoint := "/health-check"

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, endpoint, nil)

	expectedResponse := map[string]any{
		"version":      version,
		"build_date":   buildDate,
		"service_name": "api-product-catalog",
	}

	// When
	app.serveHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, expectedResponse, jsonResponse(t, w.Body.Bytes()))
}

// Test Create endpoint
func TestServer_Create(t *testing.T) {
	dataTable := []struct {
		name            string
		body            io.Reader
		productsService service.Service
		expectedCode    int
	}{
		{
			name:            "failed to create product, incorrect request body format",
			body:            mockRequest(map[string]string{"qty": "quantity"}),
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
		},
		{
			name:            "failed to create product, name cannot be empty",
			body:            mockRequest(model.Product{}),
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
		},
		{
			name: "failed to create product, sku cannot be empty",
			body: mockRequest(model.Product{
				Name: "Name1",
			}),
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
		},
		{
			name: "failed to create product, price value is invalid",
			body: mockRequest(model.Product{
				Name: "Name1",
				Sku:  "Sku1",
			}),
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
		},
		{
			name: "failed to create product, error on service",
			body: mockRequest(model.Product{
				Name:  "Name1",
				Sku:   "Sku1",
				Price: 500,
			}),
			expectedCode: http.StatusInternalServerError,
			productsService: &mockService{
				err: errors.New("error on service"),
			},
		},
		{
			name: "failed to create product, product sku already exist",
			body: mockRequest(model.Product{
				Name:  "Name1",
				Sku:   "Sku1",
				Price: 500,
			}),
			expectedCode: http.StatusBadRequest,
			productsService: &mockService{
				err: internalErrors.ErrProductSKUAlreadyExist,
			},
		},
		{
			name: "successfully create product",
			body: mockRequest(model.Product{
				Name:  "Name1",
				Sku:   "Sku1",
				Price: 500,
			}),
			expectedCode:    http.StatusCreated,
			productsService: &mockService{},
		},
	}

	for _, dt := range dataTable {
		// Given
		app := New(
			dt.productsService,
			mux.NewRouter(),
			"",
			"",
			"")

		endpoint := "/v1"

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, endpoint, dt.body)

		// When
		app.serveHTTP(w, req)

		// Then
		assert.Equal(t, dt.expectedCode, w.Code)
	}
}

// Test Get By ID endpoint
func TestServer_GetByID(t *testing.T) {
	dataTable := []struct {
		name            string
		productsService service.Service
		expectedCode    int
	}{
		{
			name:         "failed to get product, error on service",
			expectedCode: http.StatusInternalServerError,
			productsService: &mockService{
				err: errors.New("error on service"),
			},
		},
		{
			name:         "failed to get product, product not found",
			expectedCode: http.StatusNotFound,
			productsService: &mockService{
				err: internalErrors.ErrProductNotFound,
			},
		},
		{
			name:            "successfully get product",
			expectedCode:    http.StatusOK,
			productsService: &mockService{},
		},
	}

	for _, dt := range dataTable {
		// Given
		app := New(
			dt.productsService,
			mux.NewRouter(),
			"",
			"",
			"")

		endpoint := "/v1/1/"

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		// When
		app.serveHTTP(w, req)

		// Then
		assert.Equal(t, dt.expectedCode, w.Code)
	}
}

// Test Get By SKU endpoint
func TestServer_GetBySKU(t *testing.T) {
	dataTable := []struct {
		name            string
		productsService service.Service
		expectedCode    int
	}{
		{
			name:         "failed to get product, error on service",
			expectedCode: http.StatusInternalServerError,
			productsService: &mockService{
				err: errors.New("error on service"),
			},
		},
		{
			name:         "failed to get product, product not found",
			expectedCode: http.StatusNotFound,
			productsService: &mockService{
				err: internalErrors.ErrProductNotFound,
			},
		},
		{
			name:            "successfully get product",
			expectedCode:    http.StatusOK,
			productsService: &mockService{},
		},
	}

	for _, dt := range dataTable {
		// Given
		app := New(
			dt.productsService,
			mux.NewRouter(),
			"",
			"",
			"")

		endpoint := "/v1/sku/1/"

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		// When
		app.serveHTTP(w, req)

		// Then
		assert.Equal(t, dt.expectedCode, w.Code)
	}
}

// Test Delete endpoint
func TestServer_Delete(t *testing.T) {
	dataTable := []struct {
		name            string
		productsService service.Service
		expectedCode    int
	}{
		{
			name:         "failed to delete product, error on service",
			expectedCode: http.StatusInternalServerError,
			productsService: &mockService{
				err: errors.New("error on service"),
			},
		},
		{
			name:            "successfully remove product",
			expectedCode:    http.StatusNoContent,
			productsService: &mockService{},
		},
	}

	for _, dt := range dataTable {
		// Given
		app := New(
			dt.productsService,
			mux.NewRouter(),
			"",
			"",
			"")

		endpoint := "/v1/1/"

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodDelete, endpoint, nil)

		// When
		app.serveHTTP(w, req)

		// Then
		assert.Equal(t, dt.expectedCode, w.Code)
	}
}

// Test Update endpoint
func TestServer_Update(t *testing.T) {
	dataTable := []struct {
		name            string
		body            io.Reader
		productsService service.Service
		expectedCode    int
	}{
		{
			name:            "failed to update product, incorrect request body format",
			body:            mockRequest(map[string]string{"qty": "quantity"}),
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
		},
		{
			name: "successfully update product",
			body: mockRequest(model.UpdateRequest{
				Qty: 100,
			}),
			expectedCode:    http.StatusCreated,
			productsService: &mockService{},
		},
	}

	for _, dt := range dataTable {
		// Given
		app := New(
			dt.productsService,
			mux.NewRouter(),
			"",
			"",
			"")

		endpoint := "/v1/1/"

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPut, endpoint, dt.body)

		// When
		app.serveHTTP(w, req)

		// Then
		assert.Equal(t, dt.expectedCode, w.Code)
	}
}

// Test Search endpoint
func TestServer_Search(t *testing.T) {
	dataTable := []struct {
		name            string
		productsService service.Service
		expectedCode    int
		limit           string
		offset          string
	}{
		{
			name:            "failed to search product, incorrect limit format",
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
			limit:           "fake",
			offset:          "0",
		},
		{
			name:            "failed to search product, incorrect limit format",
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
			limit:           "-10",
			offset:          "0",
		},
		{
			name:            "failed to search product, incorrect offset format",
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
			limit:           "10",
			offset:          "fake",
		},
		{
			name:            "failed to search product, incorrect offset format",
			expectedCode:    http.StatusBadRequest,
			productsService: &mockService{},
			limit:           "10",
			offset:          "-20",
		},
		{
			name:         "failed to search product, error on service",
			expectedCode: http.StatusInternalServerError,
			productsService: &mockService{
				err: errors.New("error on service"),
			},
			limit:  "10",
			offset: "0",
		},
		{
			name:            "successfully search products",
			expectedCode:    http.StatusOK,
			productsService: &mockService{},
			limit:           "10",
			offset:          "0",
		},
	}

	for _, dt := range dataTable {
		// Given
		app := New(
			dt.productsService,
			mux.NewRouter(),
			"",
			"",
			"")

		endpoint := fmt.Sprintf("/v1?limit=%s&offset=%s", dt.limit, dt.offset)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		// When
		app.serveHTTP(w, req)

		// Then
		assert.Equal(t, dt.expectedCode, w.Code)
	}
}

type mockService struct {
	err error
}

func (m *mockService) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	return &model.Product{}, m.err
}

func (m *mockService) GetByID(ctx context.Context, id string) (*model.Product, error) {
	return &model.Product{}, m.err
}

func (m *mockService) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	return &model.Product{}, m.err
}

func (m *mockService) Delete(ctx context.Context, id string) error {
	return m.err
}

func (m *mockService) Update(ctx context.Context, request *model.Update) (*model.Product, error) {
	return &model.Product{}, m.err
}

func (m *mockService) Search(ctx context.Context, request model.SearchRequest) (*model.SearchResponse, error) {
	return &model.SearchResponse{}, m.err
}

func jsonResponse(t *testing.T, b []byte) map[string]any {
	var res map[string]any

	err := json.Unmarshal(b, &res)

	require.NoError(t, err)

	return res
}

func mockRequest(request any) io.Reader {
	d, _ := json.Marshal(request)

	return bytes.NewReader(d)
}
