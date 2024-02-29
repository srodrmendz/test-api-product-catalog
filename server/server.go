package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srodrmendz/api-product-catalog/docs"
	"github.com/srodrmendz/api-product-catalog/service"
	swagger "github.com/swaggo/http-swagger/v2"
)

// Initialize all dependencies
func New(productsService service.Service, router *mux.Router, path string, version string, buildDate string) *App {
	app := &App{
		Services: services{
			ProductsService: productsService,
		},
		Config: config{
			Version:   version,
			BuildDate: buildDate,
		},
		Router: router,
	}

	// Initializing panic recovery middleware
	router.Use(app.panicRecoveryMiddleware)

	subrouter := app.Router.PathPrefix(path).Subrouter()

	// Initializing healthcheck route
	subrouter.HandleFunc("/health-check", app.healthCheck).Methods(http.MethodGet)

	// Initializing create product route
	subrouter.HandleFunc("/v1", app.create).Methods(http.MethodPost)

	// Initializing search products route
	subrouter.HandleFunc("/v1", app.search).Methods(http.MethodGet)

	// Initializing get product by id
	subrouter.HandleFunc("/v1/{id}/", app.getByID).Methods(http.MethodGet)

	// Initializing delete product
	subrouter.HandleFunc("/v1/{id}/", app.delete).Methods(http.MethodDelete)

	// Initializing update product
	subrouter.HandleFunc("/v1/{id}/", app.update).Methods(http.MethodPut)

	// Initializing get product by sku
	subrouter.HandleFunc("/v1/sku/{sku}/", app.getBySKU).Methods(http.MethodGet)

	// Initializing swagger route
	router.PathPrefix(fmt.Sprintf("%s/docs", path)).Handler(swagger.WrapHandler)

	// Swagger doc
	docs.SwaggerInfo.Title = "Product Catalog API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = path
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	return app
}
