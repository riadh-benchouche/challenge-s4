package errors_test

import (
	"backend/controllers"
	"backend/database"
	"backend/models"
	"backend/services"
	"backend/tests/test_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateParticipation_Errors(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	participationService := services.NewParticipationService(database.CurrentDatabase)
	controller := controllers.NewParticipationController(participationService)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/participations", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.CreateParticipation(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("InvalidData", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		if err := database.CurrentDatabase.Create(user).Error; err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		var createdUser models.User
		if err := database.CurrentDatabase.First(&createdUser, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/participations", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", createdUser)

		err := controller.CreateParticipation(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("EventNotFound", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		if err := database.CurrentDatabase.Create(user).Error; err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		var createdUser models.User
		if err := database.CurrentDatabase.First(&createdUser, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		requestBody := `{
			"event_id": "nonexistent",
			"is_attending": true
		}`

		req := httptest.NewRequest(http.MethodPost, "/participations", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", createdUser)

		err := controller.CreateParticipation(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
