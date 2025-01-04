// tests/test_services/user_service_test.go
package services_test

import (
	"backend/enums"
	"backend/services"
	"backend/tests/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_AddUser(t *testing.T) {
	err := test_utils.SetupTestDB()
	assert.NoError(t, err)

	service := services.NewUserService()

	t.Run("Success", func(t *testing.T) {
		plainPassword := "TestPassword123!"
		user := test_utils.GetValidUser(enums.UserRole)
		user.PlainPassword = &plainPassword

		result, err := service.AddUser(user)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, user.Name, result.Name)
		assert.Equal(t, user.Email, result.Email)
		assert.NotEmpty(t, result.Password)
		assert.Nil(t, result.PlainPassword)
	})

	t.Run("MissingPassword", func(t *testing.T) {
		user := test_utils.GetValidUser(enums.UserRole)
		user.PlainPassword = nil

		result, err := service.AddUser(user)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		plainPassword := "TestPassword123!"

		// Create first user
		user1 := test_utils.GetValidUser(enums.UserRole)
		user1.PlainPassword = &plainPassword
		result1, err := service.AddUser(user1)
		assert.NoError(t, err)
		assert.NotNil(t, result1)

		// Try to create second user with same email
		user2 := test_utils.GetValidUser(enums.UserRole)
		user2.Email = user1.Email
		user2.PlainPassword = &plainPassword
		result2, err := service.AddUser(user2)
		assert.Error(t, err)
		assert.Nil(t, result2)
	})
}
