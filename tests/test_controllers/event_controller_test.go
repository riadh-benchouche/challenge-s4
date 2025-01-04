package controllers_test

import (
	"backend/controllers"
	"backend/database"
	"backend/models"
	"backend/routers"
	"backend/services"
	"backend/tests/test_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEventRouter_Unauthorized(t *testing.T) {
	// Setup
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	e := echo.New()
	service := services.NewEventService(database.CurrentDatabase)
	controller := controllers.NewEventController(service)
	routers.SetupEventRoutes(e, controller)

	tests := []struct {
		method, path string
	}{
		{http.MethodGet, "/api/events"},
		{http.MethodPost, "/api/events"},
		{http.MethodGet, "/api/events/1"},
		{http.MethodPut, "/api/events/1"},
		{http.MethodDelete, "/api/events/1"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	}
}

func TestEventRouter_CreateEvent(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	service := services.NewEventService(database.CurrentDatabase)
	controller := controllers.NewEventController(service)
	routers.SetupEventRoutes(e, controller)

	// Créer une association pour le test
	association := test_utils.GetValidAssociation()
	database.CurrentDatabase.Create(&association)

	eventJSON := `{
        "name": "Test Event",
        "description": "Test Description",
        "location": "Test Location",
        "association_id": "` + association.ID + `",
        "date": "2024-12-31T12:00:00Z"
    }`

	req := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(eventJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestEventRouter_GetEvent(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	service := services.NewEventService(database.CurrentDatabase)
	controller := controllers.NewEventController(service)
	routers.SetupEventRoutes(e, controller)

	// Créer un événement pour le test
	event := models.Event{
		Name:          "Test Event",
		Description:   "Test Description",
		Location:      "Test Location",
		AssociationID: test_utils.GetValidAssociation().ID,
	}
	database.CurrentDatabase.Create(&event)

	req := httptest.NewRequest(http.MethodGet, "/api/events/"+event.ID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestEventRouter_GetEvents(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	service := services.NewEventService(database.CurrentDatabase)
	controller := controllers.NewEventController(service)
	routers.SetupEventRoutes(e, controller)

	// Créer quelques événements pour le test
	events := []models.Event{
		{
			Name:          "Event 1",
			Description:   "Description 1",
			Location:      "Location 1",
			AssociationID: test_utils.GetValidAssociation().ID,
		},
		{
			Name:          "Event 2",
			Description:   "Description 2",
			Location:      "Location 2",
			AssociationID: test_utils.GetValidAssociation().ID,
		},
	}
	for _, event := range events {
		database.CurrentDatabase.Create(&event)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/events", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestEventRouter_UpdateEvent(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	service := services.NewEventService(database.CurrentDatabase)
	controller := controllers.NewEventController(service)
	routers.SetupEventRoutes(e, controller)

	// Créer un événement à mettre à jour
	event := models.Event{
		Name:          "Original Event",
		Description:   "Original Description",
		Location:      "Original Location",
		AssociationID: test_utils.GetValidAssociation().ID,
	}
	database.CurrentDatabase.Create(&event)

	updateJSON := `{
        "name": "Updated Event",
        "description": "Updated Description",
        "location": "Updated Location"
    }`

	req := httptest.NewRequest(http.MethodPut, "/api/events/"+event.ID, strings.NewReader(updateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}
