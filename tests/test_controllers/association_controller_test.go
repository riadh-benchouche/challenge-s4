package controllers_test

import (
	"backend/controllers"
	"backend/database"
	"backend/tests/test_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// tests/test_controllers/association_controller_test.go
func TestCreateAssociation_Integration(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	controller := controllers.NewAssociationController()

	t.Run("Success", func(t *testing.T) {
		requestBody := `{
            "name": "Test Association",
            "description": "Test Description"
        }`

		req := httptest.NewRequest(http.MethodPost, "/associations", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Créer un utilisateur authentifié
		user := test_utils.GetAuthenticatedUser()
		if err := database.CurrentDatabase.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		c.Set("user", user)

		err := controller.CreateAssociation(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
