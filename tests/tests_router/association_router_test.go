package routers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/routers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAssociationRouter_SetupRoutes(t *testing.T) {
	e := echo.New()

	// Initialisation du routeur
	r := routers.AssociationRouter{}
	r.SetupRoutes(e)

	// Table de tests
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GetAllAssociations", http.MethodGet, "/associations", http.StatusUnauthorized},
		{"GetAssociationById", http.MethodGet, "/associations/1", http.StatusUnauthorized},
		{"CreateAssociation", http.MethodPost, "/associations", http.StatusUnauthorized},
		{"UploadProfileImage", http.MethodPost, "/associations/1/upload-image", http.StatusUnauthorized},
		{"GetNextEvent", http.MethodGet, "/associations/1/next-event", http.StatusUnauthorized},
		{"GetAssociationEvents", http.MethodGet, "/associations/1/events", http.StatusUnauthorized},
		{"JoinAssociation", http.MethodPost, "/associations/join/code123", http.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Crée une requête HTTP sans authentification
			req := httptest.NewRequest(test.method, test.path, nil)
			rec := httptest.NewRecorder()

			// Simule la requête HTTP
			e.ServeHTTP(rec, req)

			// Vérifie le statut HTTP attendu (401 Unauthorized à cause du middleware)
			assert.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}
