package routers_test

import (
	"backend/routers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRouter_SetupRoutes(t *testing.T) {
	e := echo.New()

	r := routers.CategoryRouter{}
	r.SetupRoutes(e)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"CreateCategory", http.MethodPost, "/categories", http.StatusUnauthorized},
		{"GetCategories", http.MethodGet, "/categories", http.StatusUnauthorized},
		{"GetCategoryByID", http.MethodGet, "/categories/1", http.StatusUnauthorized},
		{"UpdateCategory", http.MethodPut, "/categories/1", http.StatusUnauthorized},
		{"DeleteCategory", http.MethodDelete, "/categories/1", http.StatusUnauthorized},
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
