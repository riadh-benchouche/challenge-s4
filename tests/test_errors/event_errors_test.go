package errors_test

import (
	"backend/controllers"
	"backend/database"
	"backend/models"
	"backend/tests/test_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent_Errors(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	controller := controllers.NewEventController()

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.CreateEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidData", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		if err := database.CurrentDatabase.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		var createdUser models.User
		if err := database.CurrentDatabase.First(&createdUser, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader("invalid json"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", createdUser)

		err := controller.CreateEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("UserNotInAssociation", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		if err := database.CurrentDatabase.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		var createdUser models.User
		if err := database.CurrentDatabase.First(&createdUser, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		association := test_utils.GetValidAssociation()
		requestBody := `{
            "name": "Test Event",
            "description": "Test Description",
            "association_id": "` + association.ID + `"
        }`

		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", createdUser)

		err := controller.CreateEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
