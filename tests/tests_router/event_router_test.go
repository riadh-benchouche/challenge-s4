package routers_test

import (
	"backend/routers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEventRouter_SetupRoutes(t *testing.T) {
	e := echo.New()

	r := routers.EventRouter{}
	r.SetupRoutes(e)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GetUserEventParticipation", http.MethodGet, "/events/1/participation", http.StatusUnauthorized},
		{"GetEventParticipations", http.MethodGet, "/events/1/participations", http.StatusUnauthorized},
		{"CreateEvent", http.MethodPost, "/events", http.StatusUnauthorized},
		{"GetEvents", http.MethodGet, "/events", http.StatusUnauthorized},
		{"GetEventById", http.MethodGet, "/events/1", http.StatusUnauthorized},
		{"UpdateEvent", http.MethodPut, "/events/1", http.StatusUnauthorized},
		{"DeleteEvent", http.MethodDelete, "/events/1", http.StatusUnauthorized},
		{"ChangeAttend", http.MethodPost, "/events/1/user-event-participation", http.StatusUnauthorized},
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
