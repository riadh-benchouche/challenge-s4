package services_test

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"backend/tests/test_utils"
	"backend/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddEvent_Success(t *testing.T) {
	// Setup
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	// Créer d'abord l'utilisateur
	user := test_utils.GetAuthenticatedUser()
	err = database.CurrentDatabase.Create(&user).Error
	assert.NoError(t, err)

	// Créer la catégorie
	category := test_utils.GetValidCategory()
	err = database.CurrentDatabase.Create(&category).Error
	assert.NoError(t, err)

	// Créer l'association avec le bon ownerID
	association := test_utils.GetValidAssociation()
	association.OwnerID = user.ID // Définir le propriétaire correctement
	association.Owner = user      // Définir la relation
	err = database.CurrentDatabase.Create(&association).Error
	assert.NoError(t, err)

	// Créer l'événement
	event := &models.Event{
		ID:            utils.GenerateULID(),
		Name:          "Test Event",
		Description:   "Test Description",
		Location:      "Test Location",
		Date:          time.Now(),
		CategoryID:    category.ID,
		AssociationID: association.ID,
	}

	// Test
	service := services.NewEventService(database.CurrentDatabase)
	err = service.Create(event)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, event.ID)
}

func TestGetEvents_Success(t *testing.T) {
	// Setup
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	// Créer l'utilisateur
	user := test_utils.GetAuthenticatedUser()
	err = database.CurrentDatabase.Create(&user).Error
	assert.NoError(t, err)

	// Créer la catégorie
	category := test_utils.GetValidCategory()
	err = database.CurrentDatabase.Create(&category).Error
	assert.NoError(t, err)

	// Créer l'association
	association := test_utils.GetValidAssociation()
	association.OwnerID = user.ID
	association.Owner = user
	err = database.CurrentDatabase.Create(&association).Error
	assert.NoError(t, err)

	// Créer les événements avec les bonnes relations
	events := []models.Event{
		{
			ID:            utils.GenerateULID(),
			Name:          "Event 1",
			Description:   "Description 1",
			Date:          time.Now(),
			Location:      "Location 1",
			CategoryID:    category.ID,
			AssociationID: association.ID,
		},
		{
			ID:            utils.GenerateULID(),
			Name:          "Event 2",
			Description:   "Description 2",
			Date:          time.Now(),
			Location:      "Location 2",
			CategoryID:    category.ID,
			AssociationID: association.ID,
		},
	}

	for _, event := range events {
		err = database.CurrentDatabase.Create(&event).Error
		assert.NoError(t, err)
	}

	// Test
	service := services.NewEventService(database.CurrentDatabase)
	results, err := service.GetAll()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestAddEvent_ValidationError(t *testing.T) {
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	service := services.NewEventService(database.CurrentDatabase)

	event := &models.Event{
		ID:   utils.GenerateULID(),
		Name: "Test Event",
		// Laisser AssociationID vide pour tester la validation
	}

	err = service.Create(event)
	assert.Error(t, err) // Doit retourner une erreur car AssociationID est requis
}

func TestGetEventByID_NotFound(t *testing.T) {
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	service := services.NewEventService(database.CurrentDatabase)

	// Tenter de récupérer un événement inexistant
	_, err = service.GetByID("non-existent-id")
	assert.Error(t, err)
}

func TestDeleteEvent_Success(t *testing.T) {
	// Setup
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	// Créer l'utilisateur
	user := test_utils.GetAuthenticatedUser()
	err = database.CurrentDatabase.Create(&user).Error
	assert.NoError(t, err)

	// Créer la catégorie
	category := test_utils.GetValidCategory()
	err = database.CurrentDatabase.Create(&category).Error
	assert.NoError(t, err)

	// Créer l'association
	association := test_utils.GetValidAssociation()
	association.OwnerID = user.ID
	association.Owner = user
	err = database.CurrentDatabase.Create(&association).Error
	assert.NoError(t, err)

	// Créer l'événement
	event := &models.Event{
		ID:            utils.GenerateULID(),
		Name:          "Test Event",
		Description:   "Test Description",
		Location:      "Test Location",
		Date:          time.Now(),
		CategoryID:    category.ID,
		AssociationID: association.ID,
	}
	err = database.CurrentDatabase.Create(event).Error
	assert.NoError(t, err)

	service := services.NewEventService(database.CurrentDatabase)

	// Test
	err = service.Delete(event.ID)
	assert.NoError(t, err)

	// Vérifier la suppression
	_, err = service.GetByID(event.ID)
	assert.Error(t, err)
}
