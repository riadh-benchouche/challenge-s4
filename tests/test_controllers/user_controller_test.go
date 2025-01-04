// tests/test_controllers/user_controller_test.go
package controllers_test

import (
	"backend/controllers"
	"backend/database"
	"backend/enums"
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
		// Create admin user first
		plainPassword := "TestPassword123!"
		admin := test_utils.GetValidUser(enums.AdminRole)
		admin.PlainPassword = &plainPassword
		admin.Password = "hashedpassword"
		if err := database.CurrentDatabase.Create(&admin).Error; err != nil {
			t.Fatalf("Failed to create admin user: %v", err)
		}

		// Préparer le corps de la requête
		requestBody := `{
			"name": "Test User",
			"email": "new.test@example.com",
			"password": "TestPassword123!",
			"role": "user"
		}`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", admin)

		// Log pour déboguer
		t.Logf("Request Body: %s", requestBody)

		err := controller.CreateUser(c)
		if err != nil {
			t.Logf("Controller error: %v", err)
		}

		// Log pour déboguer
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
		if err := database.CurrentDatabase.Create(&admin).Error; err != nil {
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
