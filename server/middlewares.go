package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/srodrmendz/api-auth/utils"
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

// Middleware used for authenticated routes
func (a *App) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tkHeader := r.Header.Get("Authorization")

		spl := strings.Split(tkHeader, "Bearer ")

		if len(spl) != bearerTokenHeaderLength {
			utils.ErrJSON(w, http.StatusUnauthorized, errors.New("unauthorized"))

			return
		}

		if _, err := utils.GetClaimsFromToken(spl[1], a.Config.secretKey); err != nil {
			utils.ErrJSON(w, http.StatusUnauthorized, errors.New("unauthorized"))

			return
		}

		next.ServeHTTP(w, r)

		return
	}
}
