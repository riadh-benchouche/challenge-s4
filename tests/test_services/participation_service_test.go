package services_test

import (
	"backend/database"
	"backend/services"
	"backend/tests/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParticipationService(t *testing.T) {
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	service := services.NewParticipationService(database.CurrentDatabase)

	t.Run("CreateParticipation_Success", func(t *testing.T) {

		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)
		err = database.CurrentDatabase.Create(&event).Error
		assert.NoError(t, err)

		participation := test_utils.GetValidParticipation(user.ID, event.ID)

		err = service.Create(&participation)
		assert.NoError(t, err)
		assert.NotEmpty(t, participation.ID)
	})

	t.Run("GetParticipation_Success", func(t *testing.T) {

		user := test_utils.GetAuthenticatedUser()
		err := database.CurrentDatabase.Create(&user).Error
		assert.NoError(t, err)

		association := test_utils.GetValidAssociation()
		association.OwnerID = user.ID
		err = database.CurrentDatabase.Create(&association).Error
		assert.NoError(t, err)

		event := test_utils.GetValidEvent(association.ID)
		err = database.CurrentDatabase.Create(&event).Error
		assert.NoError(t, err)

		participation := test_utils.GetValidParticipation(user.ID, event.ID)
		err = service.Create(&participation)
		assert.NoError(t, err)

		retrieved, err := service.GetByID(participation.ID)
		assert.NoError(t, err)
		assert.Equal(t, participation.UserID, retrieved.UserID)
		assert.Equal(t, participation.EventID, retrieved.EventID)
	})
}
