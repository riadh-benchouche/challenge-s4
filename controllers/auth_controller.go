package controllers

import (
	"backend/database"
	coreErrors "backend/errors"
	"backend/models"
	"backend/requests"
	"backend/services"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (c *AuthController) Login(ctx echo.Context) error {
	var jsonBody requests.LoginRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(jsonBody)
	if err != nil {
		validationErrors := utils.GetValidationErrors(err.(validator.ValidationErrors), jsonBody)
		return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
	}

	result, err := c.authService.Login(jsonBody.Email, jsonBody.Password)
	if err != nil {
		if errors.Is(err, coreErrors.ErrInvalidCredentials) {
			return ctx.String(http.StatusUnauthorized, "Invalid credentials")
		}
		if errors.Is(err, coreErrors.ErrUserNotActive) {
			return ctx.String(http.StatusForbidden, "User is not active")
		}
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) Register(ctx echo.Context) error {
	var jsonBody requests.RegisterRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(jsonBody)
	if err != nil {
		validationErrors := utils.GetValidationErrors(err.(validator.ValidationErrors), jsonBody)
		return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
	}

	var existingUser models.User
	database.CurrentDatabase.Where("email = ?", jsonBody.Email).First(&existingUser)
	if existingUser.ID != "" {
		return ctx.String(http.StatusConflict, "Email already used")
	}

	result, err := c.authService.Register(jsonBody)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, result)
}
