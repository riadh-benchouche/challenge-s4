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

func TestCreateAssociation_Integration(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	controller := controllers.NewAssociationController()

	t.Run("Success", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(user).Error
		assert.NoError(t, err, "Failed to create test user")

		t.Logf("Created user: %+v", user)

		requestBody := `{
			"name": "Test Association",
			"description": "Test Description"
		}`

		req := httptest.NewRequest(http.MethodPost, "/associations", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", user)

		err = controller.CreateAssociation(c)
		if err != nil {
			t.Logf("Controller error: %v", err)
		}
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
