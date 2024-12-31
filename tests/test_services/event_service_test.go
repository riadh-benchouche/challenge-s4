package services_test

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"backend/tests/test_utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddEvent_Success(t *testing.T) {
	err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	service := services.NewEventService(database.CurrentDatabase)

	// Créer d'abord une association pour le test
	association := test_utils.GetValidAssociation()
	database.CurrentDatabase.Create(&association)

	// Créer l'événement
	event := &models.Event{
		Name:          "Test Event",
		Description:   "Test Description",
		Location:      "Test Location",
		AssociationID: association.ID,
		Date:          time.Now(),
	}

	err = service.Create(event)

	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.NotEmpty(t, event.ID)
	assert.Equal(t, "Test Event", event.Name)
	assert.Equal(t, "Test Description", event.Description)
	assert.Equal(t, "Test Location", event.Location)
	assert.Equal(t, association.ID, event.AssociationID)
}

func TestAddEvent_ValidationError(t *testing.T) {
	err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	service := services.NewEventService(database.CurrentDatabase)

	// Créer un événement invalide (sans AssociationID requis)
	event := &models.Event{
		Name: "Test Event",
	}

	err = service.Create(event)
	assert.Error(t, err)
}

func TestGetEventByID_Success(t *testing.T) {
	err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	service := services.NewEventService(database.CurrentDatabase)

	// Créer un événement pour le test
	event := &models.Event{
		Name:          "Test Event",
		Description:   "Test Description",
		Location:      "Test Location",
		AssociationID: test_utils.GetValidAssociation().ID,
		Date:          time.Now(),
	}
	database.CurrentDatabase.Create(event)

	foundEvent, err := service.GetByID(event.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundEvent)
	assert.Equal(t, event.ID, foundEvent.ID)
}

func TestGetEventByID_NotFound(t *testing.T) {
	err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	service := services.NewEventService(database.CurrentDatabase)

	// Tenter de récupérer un événement inexistant
	_, err = service.GetByID("non-existent-id")
	assert.Error(t, err)
}

func TestGetEvents_Success(t *testing.T) {
	err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	service := services.NewEventService(database.CurrentDatabase)

	// Créer quelques événements pour le test
	association := test_utils.GetValidAssociation()
	database.CurrentDatabase.Create(&association)

	events := []models.Event{
		{
			Name:          "Event 1",
			Description:   "Description 1",
			Location:      "Location 1",
			AssociationID: association.ID,
			Date:          time.Now(),
		},
		{
			Name:          "Event 2",
			Description:   "Description 2",
			Location:      "Location 2",
			AssociationID: association.ID,
			Date:          time.Now(),
		},
	}

	for _, event := range events {
		database.CurrentDatabase.Create(&event)
	}

	allEvents, err := service.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, allEvents)
	assert.GreaterOrEqual(t, len(allEvents), 2)
}

func TestDeleteEvent_Success(t *testing.T) {
	err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	service := services.NewEventService(database.CurrentDatabase)

	// Créer un événement pour le test
	event := &models.Event{
		Name:          "Test Event",
		Description:   "Test Description",
		Location:      "Test Location",
		AssociationID: test_utils.GetValidAssociation().ID,
		Date:          time.Now(),
	}
	database.CurrentDatabase.Create(event)

	// Supprimer l'événement
	err = service.Delete(event.ID)
	assert.NoError(t, err)

	// Vérifier que l'événement a bien été supprimé
	_, err = service.GetByID(event.ID)
	assert.Error(t, err)
}
