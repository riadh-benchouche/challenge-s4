package controllers

import (
	coreErrors "backend/errors"
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
	var jsonBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=72"`
	}
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
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, result)
}
