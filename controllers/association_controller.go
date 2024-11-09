package controllers

import (
	"backend/database"
	coreErrors "backend/errors"
	"backend/models"
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

type AssociationController struct {
	AssociationService *services.AssociationService
	UserService        *services.UserService
}

func NewAssociationController() *AssociationController {
	return &AssociationController{
		AssociationService: services.NewAssociationService(),
		UserService:        services.NewUserService(),
	}
}

func (c *AssociationController) CreateAssociation(ctx echo.Context) error {
	var jsonBody models.Association

	if err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	jsonBody.ID = utils.GenerateULID()
	jsonBody.OwnerID = user.ID
	jsonBody.Owner = user
	jsonBody.Code = utils.GenerateAssociationCode()

	newAssociation, err := c.AssociationService.CreateAssociation(jsonBody)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			validationErrors := utils.GetValidationErrors(validationErrs, jsonBody)
			return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
		}
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	joinedAssociation, err := c.UserService.JoinAssociation(user.ID, newAssociation.ID, newAssociation.Code)
	if err != nil {
		switch {
		case errors.Is(err, coreErrors.ErrAlreadyJoined):
			return ctx.String(http.StatusConflict, err.Error())
		case errors.Is(err, coreErrors.ErrCodeDoesNotExist):
			return ctx.String(http.StatusNotFound, err.Error())
		default:
			ctx.Logger().Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(http.StatusCreated, joinedAssociation)
}

func (c *AssociationController) GetAssociationById(ctx echo.Context) error {
	id := ctx.Param("associationId")

	if _, err := ulid.Parse(id); err != nil {
		return ctx.JSON(http.StatusBadRequest, coreErrors.ErrInvalidULIDFormat)
	}

	association, err := c.AssociationService.GetAssociationById(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, coreErrors.ErrAssociationNotFound)
	}

	return ctx.JSON(http.StatusOK, association)
}

func (c *AssociationController) GetAllAssociations(ctx echo.Context) error {
	var filters []services.AssociationFilter
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

	associations, err := c.AssociationService.GetAllAssociations(pagination, filters...)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, associations)
}

func (a *AssociationController) UploadProfileImage(ctx echo.Context) error {
	associationID := ctx.Param("id")

	authUserID := ctx.Get("user").(models.User).ID

	var association models.Association
	if err := database.CurrentDatabase.First(&association, "id = ?", associationID).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, "Association not found")
	}

	if association.OwnerID != authUserID {
		return ctx.JSON(http.StatusForbidden, "You do not have permission to update this association")
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

	imageName := associationID + "_" + file.Filename
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
	err = database.CurrentDatabase.Model(&models.Association{}).
		Where("id = ?", associationID).
		Update("image_url", imagePath).Error

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de la mise à jour de l'URL de l'image")
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Image uploadée avec succès", "image_url": imagePath})
}
