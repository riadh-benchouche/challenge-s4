// tests/errors/association_errors_test.go
package errors_test

import (
	"backend/controllers"
	"backend/tests/test_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// tests/test_errors/association_errors_test.go
func TestCreateAssociation_Errors(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	controller := controllers.NewAssociationController()

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/associations", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.CreateAssociation(c)
		assert.NoError(t, err) // Le controller gère lui-même l'erreur
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidData", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/associations", strings.NewReader("invalid json"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Ajouter un utilisateur authentifié
		user := test_utils.GetAuthenticatedUser()
		c.Set("user", user)

		err := controller.CreateAssociation(c)
		assert.NoError(t, err) // Le controller gère lui-même l'erreur
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}