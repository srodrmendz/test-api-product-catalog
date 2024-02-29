package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srodrmendz/api-auth/docs"
	"github.com/srodrmendz/api-auth/service"
	swagger "github.com/swaggo/http-swagger/v2"
)

// Initialize all dependencies
func New(authService service.Service, router *mux.Router, path string, secretKey string, version string, buildDate string) *App {
	app := &App{
		Services: services{
			authService: authService,
		},
		Config: config{
			version:   version,
			buildDate: buildDate,
			secretKey: secretKey,
		},
		Router: router,
	}

	// Initializing panic recovery middleware
	router.Use(app.panicRecoveryMiddleware)

	subrouter := app.Router.PathPrefix(path).Subrouter()

	// Initializing healthcheck route
	subrouter.HandleFunc("/health-check", app.healthCheck).Methods(http.MethodGet)

	// Initializing authenticate route
	subrouter.HandleFunc("/v1/", app.authenticate).Methods(http.MethodPost)

	// Initializing protected route
	subrouter.HandleFunc("/protected", app.authMiddleware(app.protected)).Methods("GET")

	// Initializing swagger route
	router.PathPrefix(fmt.Sprintf("%s/docs", path)).Handler(swagger.WrapHandler)

	// Swagger doc
	docs.SwaggerInfo.Title = "Auth API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/auth"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	return app
}
