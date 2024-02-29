package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	internalErrors "github.com/srodrmendz/api-auth/errors"
	"github.com/srodrmendz/api-auth/model"
	"github.com/srodrmendz/api-auth/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Healthcheck endpoint
func TestServer_HealthCheck(t *testing.T) {
	// Given
	version := "v1.0.0"

	buildDate := "02/29/2024"

	app := New(
		nil,
		mux.NewRouter(),
		"",
		"",
		version,
		buildDate,
	)

	endpoint := "/health-check"

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, endpoint, nil)

	expectedResponse := map[string]any{
		"version":      version,
		"build_date":   buildDate,
		"service_name": "api-auth",
	}

	// When
	app.serveHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, expectedResponse, jsonResponse(t, w.Body.Bytes()))
}

// Test Authenticate endpoint
func TestServer_Authenticate(t *testing.T) {
	dataTable := []struct {
		name         string
		body         io.Reader
		expectedCode int
		authService  service.Service
	}{
		{
			name:         "failed to authenticate user, incorrect request format",
			body:         mockRequest(map[string]int{"email": 1, "password": 12}),
			expectedCode: http.StatusBadRequest,
			authService:  &mockAuthService{},
		},
		{
			name:         "failed to authenticate user, incorrect email format",
			body:         mockRequest(model.AuthRequest{Email: "fake", Password: ""}),
			expectedCode: http.StatusBadRequest,
			authService:  &mockAuthService{},
		},
		{
			name:         "failed to authenticate user, incorrect credentials",
			body:         mockRequest(model.AuthRequest{Email: "fake@gmail.com", Password: ""}),
			expectedCode: http.StatusUnauthorized,
			authService: &mockAuthService{
				err: internalErrors.ErrUserNotFound,
			},
		},
		{
			name:         "successfully authenticate user",
			body:         mockRequest(model.AuthRequest{Email: "fake@gmail.com", Password: ""}),
			expectedCode: http.StatusOK,
			authService:  &mockAuthService{},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			app := New(
				dt.authService,
				mux.NewRouter(),
				"",
				"",
				"",
				"",
			)

			endpoint := "/v1/"

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, endpoint, dt.body)

			// When
			app.serveHTTP(w, req)

			// Then
			assert.Equal(t, dt.expectedCode, w.Code)
		})
	}
}

// Test Protected endpoint
func TestServer_Protected(t *testing.T) {
	dataTable := []struct {
		name             string
		header           http.Header
		expectedCode     int
		expectedResponse map[string]any
	}{
		{
			name:             "failed to access method, token is not sended on header",
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]any{"error": "unauthorized"},
		},
		{
			name: "failed to access method, token is not sended on bearer format",
			header: http.Header{
				"Authorization": []string{"Token"},
			},
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]any{"error": "unauthorized"},
		},
		{
			name: "failed to access method, bearer token incorrect format",
			header: http.Header{
				"Authorization": []string{"Bearer Token Token"},
			},
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]any{"error": "unauthorized"},
		},
		{
			name: "successfully authenticate on method",
			header: http.Header{
				"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImZha2VAZW1haWwuY29tIiwidXNlcm5hbWUiOiJmYWtlIiwiZXhwIjo1MzA5MjI0MDUyLCJqdGkiOiJiMzNkYTU2OS04YmI3LTRlMmEtYTUyMy0zNDNiNjg5M2U5NWIiLCJpYXQiOjE3MDkyMjc2NTIsImlzcyI6InRlc3RfYXBwIn0._e3Alt4Pngmeq33uxRs7hRf3V7Uu1Ero5UHD815uSu0"},
			},
			expectedCode:     http.StatusOK,
			expectedResponse: map[string]any{"status": "OK"},
		},
	}

	for _, dt := range dataTable {
		t.Run(dt.name, func(t *testing.T) {
			// Given
			app := New(
				nil,
				mux.NewRouter(),
				"",
				"test-key",
				"",
				"",
			)

			endpoint := "/protected"

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, endpoint, nil)

			req.Header = dt.header

			// When
			app.serveHTTP(w, req)

			// Then
			assert.Equal(t, dt.expectedCode, w.Code)

			assert.Equal(t, dt.expectedResponse, jsonResponse(t, w.Body.Bytes()))
		})
	}
}

type mockAuthService struct {
	err error
}

func (m *mockAuthService) Authenticate(ctx context.Context, email string, password string) (*model.AuthResponse, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &model.AuthResponse{}, nil
}

func jsonResponse(t *testing.T, b []byte) map[string]any {
	var res map[string]any

	err := json.Unmarshal(b, &res)

	require.NoError(t, err)

	return res
}

func mockRequest(request any) io.Reader {
	d, _ := json.Marshal(request)

	return bytes.NewReader(d)
}
