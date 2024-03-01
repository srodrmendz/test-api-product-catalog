package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/srodrmendz/api-product-catalog/conf"
	"github.com/srodrmendz/api-product-catalog/repository"
	"github.com/srodrmendz/api-product-catalog/server"
	"github.com/srodrmendz/api-product-catalog/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	version   = "NO_VERSION"
	buildDate = "NO_BUILD_DATE"
)

// @contact.name Sebastian Rodriguez Mendez
// @contact.email srodmendz@gmail.com

// nolint: forbidigo
func main() {
	ctx := context.TODO()

	// Create mongo connection
	mongoClient := createMongoClient(
		ctx,
		conf.GetProps().Database.URI,
	)

	// Disconnect after server is shutdown
	defer mongoClient.Disconnect(ctx)

	// Create users repository
	repository := repository.New(
		mongoClient,
		conf.GetProps().Database.DB,
		conf.GetProps().Database.Collection,
	)

	// Create service
	service := service.New(repository)

	// Create app
	app := server.New(
		service,
		mux.NewRouter(),
		conf.GetProps().Path,
		version,
		buildDate,
	)

	fmt.Println("running on 8080 port")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	http.Handle("/", corsHandler(app.Router))

	fmt.Println("listening on 8080 port...")

	http.ListenAndServe(":8080", nil)
}

func createMongoClient(ctx context.Context, uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	return client
}
