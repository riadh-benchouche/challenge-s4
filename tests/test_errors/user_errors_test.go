// tests/test_errors/user_errors_test.go
package errors_test

import (
	"backend/controllers"
	"backend/database"
	"backend/tests/test_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Errors(t *testing.T) {
	if err := test_utils.SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	e := echo.New()
	controller := controllers.NewUserController()

	t.Run("NoAuthUser", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := controller.CreateUser(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("MissingPassword", func(t *testing.T) {
		admin := test_utils.GetAdminUser()
		if err := database.CurrentDatabase.Create(&admin).Error; err != nil {
			t.Fatalf("Failed to create admin user: %v", err)
		}

		requestBody := `{
			"name": "Test User",
			"email": "test@example.com"
		}`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", admin)

		err := controller.CreateUser(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("InvalidPasswordFormat", func(t *testing.T) {
		admin := test_utils.GetAdminUser()
		if err := database.CurrentDatabase.Create(&admin).Error; err != nil {
			t.Fatalf("Failed to create admin user: %v", err)
		}

		requestBody := `{
			"name": "Test User",
			"email": "test@example.com",
			"password": "short"
		}`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", admin)

		err := controller.CreateUser(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}
