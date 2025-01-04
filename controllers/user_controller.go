package controllers

import (
	"backend/database"
	"backend/enums"
	coreErrors "backend/errors"
	"backend/models"
	"backend/resources"
	"backend/services"
	"backend/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

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
	// 1. Bind la requête
	var jsonBody models.User
	if err := ctx.Bind(&jsonBody); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	// 2. Vérifier l'authentification
	currentUser := ctx.Get("user")
	if currentUser == nil {
		return ctx.NoContent(http.StatusForbidden)
	}

	if admin, ok := currentUser.(models.User); !ok || admin.Role != enums.AdminRole {
		return ctx.NoContent(http.StatusForbidden)
	}

	// 3. Générer l'ID avant la validation
	jsonBody.ID = utils.GenerateULID()

	// 4. Créer l'utilisateur
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

		if errors.Is(err, coreErrors.ErrInvalidPassword) {
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
				"error": "Invalid password",
			})
		}

		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, newUser)
}

func (c *UserController) GetUsers(ctx echo.Context) error {
	var filters []services.UserFilter
	params := ctx.QueryParams().Get("filters")

	if len(params) > 0 {
		err := json.Unmarshal([]byte(params), &filters)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	for _, filter := range filters {
		err := validate.Struct(filter)
		if err != nil {
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {
				validationErrors := utils.GetValidationErrors(validationErrs, filter)
				return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
			}
			ctx.Logger().Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	pagination := utils.PaginationFromContext(ctx)

	userPagination, err := c.UserService.GetUsers(pagination, filters...)
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

	userResource := resources.NewUserResource(*user)
	return ctx.JSON(http.StatusOK, userResource)
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

func (c *UserController) GetOwnerAssociations(ctx echo.Context) error {
	id := ctx.Param("id")
	if _, err := ulid.Parse(id); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, err := c.UserService.FindByID(id)

	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	associations := user.AssociationsOwned

	associationResources := make([]resources.AssociationResource, len(associations))
	for i, association := range associations {
		associationResources[i] = resources.NewAssociationResource(association)
	}

	return ctx.JSON(http.StatusOK, associationResources)
}

func (c *UserController) GetUserAssociations(ctx echo.Context) error {
	id := ctx.Param("id")
	if _, err := ulid.Parse(id); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, err := c.UserService.FindByID(id)

	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	associations := user.Associations

	associationResources := make([]resources.AssociationResource, len(associations))
	for i, association := range associations {
		associationResources[i] = resources.NewAssociationResource(association)
	}

	return ctx.JSON(http.StatusOK, associationResources)
}

func (c *UserController) JoinAssociation(ctx echo.Context) error {
	userID := ctx.Param("id")
	associationID := ctx.Param("association_id")

	if _, err := ulid.Parse(userID); err != nil {
		return ctx.JSON(http.StatusBadRequest, coreErrors.ErrInvalidULIDFormat)
	}

	if _, err := ulid.Parse(associationID); err != nil {
		return ctx.JSON(http.StatusBadRequest, coreErrors.ErrInvalidULIDFormat)
	}

	code := ctx.QueryParam("code")
	if code == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Code is required"})
	}

	if hasJoined, err := c.UserService.JoinAssociation(userID, associationID, code); err != nil {
		ctx.Logger().Error(err)

		switch {
		case errors.Is(err, coreErrors.ErrAlreadyJoined):
			return ctx.JSON(http.StatusConflict, map[string]string{"error": "User already joined"})
		case errors.Is(err, coreErrors.ErrInvalidCode):
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid association code"})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	} else if hasJoined {
		return ctx.JSON(http.StatusCreated, "User successfully joined the association")
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c *UserController) UploadProfileImage(ctx echo.Context) error {
	userID := ctx.Param("id")

	authUserId := ctx.Get("user").(models.User).ID
	if userID != authUserId {
		return ctx.JSON(http.StatusForbidden, "Vous ne pouvez pas modifier l'image de profil d'un autre utilisateur !")
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.Logger().Error("Error retrieving file: ", err)
		return ctx.JSON(http.StatusBadRequest, "Erreur lors de l'upload de l'image")
	}

	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de l'ouverture de l'image")
	}
	defer src.Close()

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.MkdirAll("public", os.ModePerm)
	}

	imageName := userID + "" + file.Filename
	imagePath := "public/" + imageName

	dst, err := os.Create(imagePath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de la création de l'image")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de l'enregistrement de l'image")
	}

	// Sauvegarde le chemin de l'image dans le champ ImageURL du modèle User
	err = database.CurrentDatabase.Model(&models.User{}).
		Where("id = ?", userID).
		Update("image_url", imagePath).Error

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de la mise à jour de l'URL de l'image")
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Image uploadée avec succès", "image_url": imagePath})
}
