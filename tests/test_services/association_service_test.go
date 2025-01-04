package services_test

import (
	"backend/database"
	"backend/services"
	"backend/tests/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Dans association_service_test.go
func TestCreateAssociation_Success(t *testing.T) {
	// Setup
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	// Créer l'utilisateur d'abord
	user := test_utils.GetAuthenticatedUser()
	err = database.CurrentDatabase.Create(&user).Error
	assert.NoError(t, err)

	// Créer l'association
	association := test_utils.GetValidAssociation()
	association.OwnerID = user.ID
	association.Owner = user

	service := services.NewAssociationService() // Sans paramètre
	result, err := service.CreateAssociation(association)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result != nil {
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, association.Name, result.Name)
	}
}
