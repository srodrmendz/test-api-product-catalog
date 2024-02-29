package server

import (
	"github.com/gorilla/mux"
	"github.com/srodrmendz/api-auth/service"
)

const bearerTokenHeaderLength = 2

type services struct {
	authService service.Service
}

type config struct {
	version   string
	buildDate string
	secretKey string
}

type App struct {
	Services services
	Config   config
	Router   *mux.Router
}
