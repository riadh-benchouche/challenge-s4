package controllers_test

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

func TestCreateParticipation_Integration(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	participationService := services.NewParticipationService(database.CurrentDatabase)
	controller := controllers.NewParticipationController(participationService)

	t.Run("Success", func(t *testing.T) {

		user := test_utils.GetAuthenticatedUser()
		if err := database.CurrentDatabase.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		var createdUser models.User
		if err := database.CurrentDatabase.First(&createdUser, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve created user: %v", err)
		}

		association := test_utils.GetValidAssociation()
		association.OwnerID = createdUser.ID
		if err := database.CurrentDatabase.Create(&association).Error; err != nil {
			t.Fatalf("Failed to create test association: %v", err)
		}

		event := test_utils.GetValidEvent(association.ID)
		if err := database.CurrentDatabase.Create(&event).Error; err != nil {
			t.Fatalf("Failed to create test event: %v", err)
		}

		requestBody := `{
            "user_id": "` + createdUser.ID + `",
            "event_id": "` + event.ID + `",
            "is_attending": true
        }`

		req := httptest.NewRequest(http.MethodPost, "/participations", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", createdUser)

		err := controller.CreateParticipation(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("Get Participations Success", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		if err := database.CurrentDatabase.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		var createdUser models.User
		if err := database.CurrentDatabase.First(&createdUser, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve created user: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/participations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", createdUser)

		err := controller.GetParticipations(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
