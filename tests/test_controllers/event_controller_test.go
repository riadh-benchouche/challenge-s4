package controllers_test

import (
	"backend/controllers"
	"backend/database"
	"backend/models"
	"backend/tests/test_utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent_Integration(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Échec configuration BD test: %v", err)
	}

	e := echo.New()
	controller := controllers.NewEventController()

	t.Run("Success", func(t *testing.T) {
		user, association := test_utils.CreateUserAndAssociation()

		var createdUser models.User
		if err := database.CurrentDatabase.First(&createdUser, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve created user: %v", err)
		}

		requestBody := fmt.Sprintf(`{
            "name": "Test Event",
            "description": "Test Description",
            "location": "Test Location",
            "association_id": "%s"
        }`, association.ID)

		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("user", createdUser)

		err := controller.CreateEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
