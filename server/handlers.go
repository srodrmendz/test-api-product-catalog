package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"

	internalErrors "github.com/srodrmendz/api-auth/errors"
	"github.com/srodrmendz/api-auth/model"
	"github.com/srodrmendz/api-auth/utils"
)

// Healthcheck godoc
// @Tags healthcheck
// @Accept  json
// @Produce  json
// @Success 200
// @Router /health-check [get]
func (a *App) healthCheck(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{
		"version":      a.Config.version,
		"build_date":   a.Config.buildDate,
		"service_name": "api-auth",
	}

	utils.DataJSON(w, http.StatusOK, response)
}

// Authenticate godoc
// @Tags authenticate
// @Description Authenticate
// @Accept  json
// @Produce  json
/// @Param request body model.AuthRequest true "Request body"
// @Success 200 {object} model.AuthResponse
// @Failure 401
// @Failure 500
// @Router /v1/ [post]
func (a *App) authenticate(w http.ResponseWriter, r *http.Request) {
	var request model.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ErrJSON(w, http.StatusBadRequest, errors.New("incorrect request body format"))

		return
	}

	if _, err := mail.ParseAddress(request.Email); err != nil {
		utils.ErrJSON(w, http.StatusBadRequest, errors.New("incorrect email format"))

		return
	}

	resp, err := a.Services.authService.Authenticate(r.Context(), request.Email, request.Password)
	if err != nil {
		if errors.Is(err, internalErrors.ErrUserNotFound) {
			utils.ErrJSON(w, http.StatusUnauthorized, err)

			return
		}

		utils.ErrJSON(w, http.StatusInternalServerError, err)

		return
	}

	utils.DataJSON(w, http.StatusOK, resp)
}

// Protected godoc
// @Tags protected
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200
// @Failure 401
// @Router /protected [get]
func (a *App) protected(w http.ResponseWriter, r *http.Request) {
	utils.DataJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}
