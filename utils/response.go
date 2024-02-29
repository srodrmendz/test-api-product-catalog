package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

// Used to return success json responses on http handlers
func DataJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	data(w, statusCode, v)
}

// Used to return errors json responses on http handlers
func ErrJSON(w http.ResponseWriter, statusCode int, err error) {
	response := errorResponse{
		Error: err.Error(),
	}

	DataJSON(w, statusCode, response)
}

func data(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		err = fmt.Errorf("encoding response %w", err)

		response := errorResponse{
			Error: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
	}
}
