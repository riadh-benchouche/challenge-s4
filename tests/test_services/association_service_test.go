package services_test

import (
	"backend/database"
	"backend/services"
	"backend/tests/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAssociation_Success(t *testing.T) {

	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	user := test_utils.GetAuthenticatedUser()
	err = database.CurrentDatabase.Create(user).Error
	assert.NoError(t, err)

	association := test_utils.GetValidAssociation()
	association.OwnerID = user.ID
	association.Owner = *user

	service := services.NewAssociationService()
	result, err := service.CreateAssociation(*association)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result != nil {
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, association.Name, result.Name)
	}
}
