// backend/tests/test_services/association_service_test.go
package services_test

import (
	"backend/services"
	"backend/tests/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAssociation_Success(t *testing.T) {
	// Setup
	err := test_utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	service := services.NewAssociationService()
	if service == nil {
		t.Fatal("Service should not be nil")
	}

	// Créer d'abord un utilisateur
	user := test_utils.GetAuthenticatedUser()
	association := test_utils.GetValidAssociation()
	association.OwnerID = user.ID
	association.Owner = user

	// Test
	result, err := service.CreateAssociation(association)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result != nil {
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, association.Name, result.Name)
	}
}
