package routers_test

import (
	"backend/routers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRouter_SetupRoutes(t *testing.T) {
	e := echo.New()

	r := routers.UserRouter{}
	r.SetupRoutes(e)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"CreateUser", http.MethodPost, "/users", http.StatusUnauthorized},     // Middleware empêche l'accès
		{"GetUsers", http.MethodGet, "/users", http.StatusUnauthorized},        // Middleware empêche l'accès
		{"UpdateUser", http.MethodPut, "/users/1", http.StatusUnauthorized},    // Middleware empêche l'accès
		{"DeleteUser", http.MethodDelete, "/users/1", http.StatusUnauthorized}, // Middleware empêche l'accès
		{"FindByID", http.MethodGet, "/users/1", http.StatusUnauthorized},      // Middleware empêche l'accès
		{"GetOwnerAssociations", http.MethodGet, "/users/1/owner-associations", http.StatusUnauthorized},
		{"GetUserAssociations", http.MethodGet, "/users/1/associations", http.StatusUnauthorized},
		{"JoinAssociation", http.MethodPost, "/users/1/associations/1", http.StatusUnauthorized},
		{"UploadProfileImage", http.MethodPost, "/users/1/upload-image", http.StatusUnauthorized},
		{"GetUserEvents", http.MethodGet, "/users/events", http.StatusUnauthorized},
		{"GetAssociationsEvents", http.MethodGet, "/users/associations/events", http.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(test.method, test.path, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}
