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

		createdUser, err := service.AddUser(*user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.NotEmpty(t, createdUser.ID)
		assert.NotEmpty(t, createdUser.Password)
		assert.Empty(t, createdUser.PlainPassword)
	})

	t.Run("MissingPassword", func(t *testing.T) {
		user := test_utils.GetValidUser(enums.UserRole)
		user.PlainPassword = nil

		_, err := service.AddUser(*user)
		assert.Error(t, err)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		plainPassword := "TestPassword123!"
		user := test_utils.GetValidUser(enums.UserRole)
		user.PlainPassword = &plainPassword

		_, err := service.AddUser(*user)
		assert.NoError(t, err)

		_, err = service.AddUser(*user)
		assert.Error(t, err)
	})
}
