package server

import (
	"encoding/json"
	"net/http"

	"github.com/srodrmendz/api-product-catalog/utils"
)

// Middleware used for not handled exceptions
func (a *App) panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				body, _ := json.Marshal(map[string]interface{}{
					"error": "internal server error",
				})

				utils.DataJSON(w, http.StatusInternalServerError, body)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *App) serveHTTP(w http.ResponseWriter, req *http.Request) {
	a.Router.ServeHTTP(w, req)
}
