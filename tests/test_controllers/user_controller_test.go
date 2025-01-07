package controllers_test

import (
	"backend/controllers"
	"backend/database"
	"backend/enums"
	"backend/models"
	"backend/tests/test_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Integration(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	controller := controllers.NewUserController()

	t.Run("Success", func(t *testing.T) {

		plainPassword := "TestPassword123!"
		admin := test_utils.GetValidUser(enums.AdminRole)
		admin.PlainPassword = &plainPassword
		admin.Password = "hashedpassword"
		admin.IsActive = true
		admin.IsConfirmed = true
		admin.Role = enums.AdminRole

		if err := database.CurrentDatabase.Create(admin).Error; err != nil {
			t.Fatalf("Failed to create admin user: %v", err)
		}

		var createdAdmin models.User
		if err := database.CurrentDatabase.First(&createdAdmin, "id = ?", admin.ID).Error; err != nil {
			t.Fatalf("Failed to retrieve admin: %v", err)
		}
		t.Logf("Created admin: %+v", createdAdmin)

		requestBody := `{
            "name": "Test User",
            "email": "new.test@example.com",
            "plain_password": "TestPassword123!",
            "role": "user"
        }`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("user", createdAdmin)

		err := controller.CreateUser(c)
		if err != nil {
			t.Logf("Controller error: %v", err)
		}

		t.Logf("Response Status: %d", rec.Code)
		t.Logf("Response Body: %s", rec.Body.String())

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("Get Users Success", func(t *testing.T) {
		plainPassword := "TestPassword123!"
		admin := test_utils.GetValidUser(enums.AdminRole)
		admin.PlainPassword = &plainPassword
		admin.Password = "hashedpassword"
		if err := database.CurrentDatabase.Create(admin).Error; err != nil {
			t.Fatalf("Failed to create admin user: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", admin)

		err := controller.GetUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
