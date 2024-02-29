package server

import (
	"github.com/gorilla/mux"
	"github.com/srodrmendz/api-product-catalog/service"
)

type services struct {
	ProductsService service.Service
}

type config struct {
	Version   string
	BuildDate string
}

type App struct {
	Services services
	Config   config
	Router   *mux.Router
}
