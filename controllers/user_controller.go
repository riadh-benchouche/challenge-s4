package controllers

import (
	"backend/enums"
	coreErrors "backend/errors"
	"backend/models"
	"backend/services"
	"backend/utils"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		UserService: services.NewUserService(),
	}
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	var jsonBody models.User
	err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	jsonBody.ID = ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

	if _, err := ulid.Parse(jsonBody.ID); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ULID format"})
	}

	currentUser := ctx.Get("user")
	if currentUser == nil || (currentUser.(models.User).Role != enums.AdminRole) {
		jsonBody.Role = enums.UserRole
	}

	newUser, err := c.UserService.AddUser(jsonBody)
	if err != nil {
		if errors.Is(err, coreErrors.ErrUserAlreadyExists) {
			return ctx.String(http.StatusConflict, err.Error())
		}

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			validationErrors := utils.GetValidationErrors(validationErrs, jsonBody)
			return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
		}

		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, newUser)
}

func (c *UserController) GetUsers(ctx echo.Context) error {
	pagination := utils.PaginationFromContext(ctx)
	search := ctx.QueryParam("search")

	userPagination, err := c.UserService.GetUsers(pagination, &search)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, userPagination)
}
func (c *UserController) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")

	if _, err := ulid.Parse(id); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	currentUser := ctx.Get("user").(models.User)

	if id == currentUser.ID {
		return ctx.String(http.StatusForbidden, "Impossible de supprimer votre propre compte.")
	}

	err := c.UserService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, coreErrors.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c *UserController) FindByID(ctx echo.Context) error {
	id := ctx.Param("id")

	_, err := ulid.Parse(id)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, err := c.UserService.FindByID(id)
	if err != nil {
		if errors.Is(err, coreErrors.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, user)
}

func (c *UserController) UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if _, err := ulid.Parse(id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ULID format"})
	}

	currentUser := ctx.Get("user").(models.User)
	if currentUser.ID != id && !enums.IsAdmin(currentUser.Role) {
		return ctx.NoContent(http.StatusForbidden)
	}

	userToUpdate, err := c.UserService.FindByID(id)
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	var jsonBody models.User
	err = json.NewDecoder(ctx.Request().Body).Decode(&jsonBody)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if jsonBody.Name != "" {
		userToUpdate.Name = jsonBody.Name
	}
	if jsonBody.Email != "" && jsonBody.Email != userToUpdate.Email {
		userToUpdate.Email = jsonBody.Email
	}

	if jsonBody.PlainPassword != nil {
		userToUpdate.PlainPassword = jsonBody.PlainPassword
	}

	if enums.IsAdmin(currentUser.Role) && currentUser.ID != id {
		userToUpdate.Role = jsonBody.Role
	} else if jsonBody.Role != "" && jsonBody.Role != userToUpdate.Role {
		return ctx.String(http.StatusForbidden, "Impossible de modifier votre rôle.")
	}

	updatedUser, err := c.UserService.UpdateUser(id, *userToUpdate)
	if err != nil {
		if errors.Is(err, coreErrors.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}

		if errors.Is(err, coreErrors.ErrUserAlreadyExists) {
			return ctx.NoContent(http.StatusConflict)
		}

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			validationErrors := utils.GetValidationErrors(validationErrs, jsonBody)
			return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
		}

		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, updatedUser)
}
