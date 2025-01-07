package services_test

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"backend/tests/test_utils"
	"backend/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventService(t *testing.T) {
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	service := services.NewEventService()

	t.Run("CreateEvent_Success", func(t *testing.T) {

		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)

		createdEvent, err := service.AddEvent(&event)
		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)
		assert.NotEmpty(t, createdEvent.ID)
		assert.Equal(t, event.Name, createdEvent.Name)
	})

	t.Run("GetEvent_Success", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)
		createdEvent, err := service.AddEvent(&event)
		assert.NoError(t, err)

		retrievedEvent, err := service.GetEventById(createdEvent.ID)
		assert.NoError(t, err)
		assert.Equal(t, event.Name, retrievedEvent.Name)
		assert.Equal(t, event.AssociationID, retrievedEvent.AssociationID)
	})

	t.Run("UpdateEvent_Success", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)
		createdEvent, err := service.AddEvent(&event)
		assert.NoError(t, err)

		createdEvent.Name = "Updated Event Name"
		createdEvent.Description = "Updated Event Description"
		createdEvent.Location = "Updated Location"

		err = service.UpdateEvent(createdEvent)
		assert.NoError(t, err)

		updatedEvent, err := service.GetEventById(createdEvent.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Event Name", updatedEvent.Name)
		assert.Equal(t, "Updated Event Description", updatedEvent.Description)
		assert.Equal(t, "Updated Location", updatedEvent.Location)
	})

	t.Run("DeleteEvent_Success", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)
		createdEvent, err := service.AddEvent(&event)
		assert.NoError(t, err)

		err = service.DeleteEvent(createdEvent.ID)
		assert.NoError(t, err)

		_, err = service.GetEventById(createdEvent.ID)
		assert.Error(t, err)
	})

	t.Run("GetEvents_WithPagination", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		for i := 0; i < 5; i++ {
			event := test_utils.GetValidEvent(association.ID)
			_, err := service.AddEvent(&event)
			assert.NoError(t, err)
		}

		pagination := utils.Pagination{
			Page:  1,
			Limit: 3,
		}
		search := "Test Event"

		result, err := service.GetEvents(pagination, &search)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 3, len(result.Rows.([]models.Event)))
	})

	t.Run("GetEventParticipations_Success", func(t *testing.T) {
		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)
		createdEvent, err := service.AddEvent(&event)
		assert.NoError(t, err)

		participation := test_utils.GetValidParticipation(user.ID, createdEvent.ID)
		err = database.CurrentDatabase.Create(&participation).Error
		assert.NoError(t, err)

		pagination := utils.Pagination{
			Page:  1,
			Limit: 10,
		}
		status := "pending"

		result, err := service.GetEventParticipations(createdEvent.ID, pagination, &status)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("ChangeUserEventAttend_Success", func(t *testing.T) {

		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)
		createdEvent, err := service.AddEvent(&event)
		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)

		participation := models.Participation{
			UserID:      user.ID,
			EventID:     createdEvent.ID,
			IsAttending: false,
		}
		err = database.CurrentDatabase.Create(&participation).Error
		assert.NoError(t, err)

		updatedParticipation, err := service.ChangeUserEventAttend(true, createdEvent.ID, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, updatedParticipation)
		assert.True(t, updatedParticipation.IsAttending)

		var checkParticipation models.Participation
		err = database.CurrentDatabase.Where("event_id = ? AND user_id = ?", createdEvent.ID, user.ID).First(&checkParticipation).Error
		assert.NoError(t, err)
		assert.True(t, checkParticipation.IsAttending)
	})
}
