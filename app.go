package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srodrmendz/api-auth/conf"
	"github.com/srodrmendz/api-auth/repository"
	"github.com/srodrmendz/api-auth/server"
	"github.com/srodrmendz/api-auth/service"
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
	service := service.New(
		repository,
		conf.GetProps().SecretKey,
	)

	// Create app
	app := server.New(
		service,
		mux.NewRouter(),
		conf.GetProps().Path,
		conf.GetProps().SecretKey,
		version,
		buildDate,
	)

	fmt.Println("running on 8080 port")

	http.ListenAndServe(":8080", app.Router)
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
