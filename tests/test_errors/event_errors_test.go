package errors_test

import (
	"backend/controllers"
	"backend/database"
	"backend/services"
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
	service := services.NewEventService(database.CurrentDatabase)
	controller := controllers.NewEventController(service)

	t.Run("InvalidData", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader("invalid json"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.Create(c)
		assert.NoError(t, err) // Le controller gère lui-même l'erreur
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("MissingAssociationID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{
            "name": "Test Event",
            "description": "Test Description"
        }`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
