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
	"gorm.io/gorm"
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

	currentUser := ctx.Get("user").(models.User)

	if currentUser.ID != id && !enums.IsAdmin(currentUser.Role) {
		return ctx.JSON(http.StatusForbidden, "Vous n'êtes pas autorisé à voir ces informations.")
	}

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

	if len(associations) == 0 {
		return ctx.JSONPretty(http.StatusNoContent, "Aucune association trouvée", " ")
	}

	return ctx.JSON(http.StatusOK, associations)
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

// UploadProfileImage TODO: fix a max image size
func (c *UserController) UploadProfileImage(ctx echo.Context) error {
	userID := ctx.Param("id")

	authUserId := ctx.Get("user").(models.User).ID
	if userID != authUserId {
		return ctx.JSON(http.StatusForbidden, "Vous ne pouvez pas modifier l'image de profil d'un autre utilisateur !")
	}

	// Récupérer le fichier depuis le formulaire
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.Logger().Error("Error retrieving file: ", err)
		return ctx.JSON(http.StatusBadRequest, "Erreur lors de l'upload de l'image")
	}

	// Vérifier la taille du fichier (exemple : 5 MB)
	const maxFileSize = 5 * 1024 * 1024 // 5 MB
	if file.Size > maxFileSize {
		return ctx.JSON(http.StatusBadRequest, "La taille du fichier dépasse la limite autorisée de 5 MB")
	}

	// Ouvrir le fichier temporairement pour vérifier son type
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de l'ouverture de l'image")
	}
	defer src.Close()

	// Lire les premiers 512 octets pour détecter le type MIME
	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de la lecture de l'image")
	}

	fileType := http.DetectContentType(buffer)

	// Types MIME acceptés (ajoutez ceux qui sont pertinents pour votre cas)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	if !allowedTypes[fileType] {
		return ctx.JSON(http.StatusBadRequest, "Type de fichier non supporté. Seules les images JPEG, PNG et GIF sont autorisées")
	}

	// Retourner au début du fichier pour le traitement ultérieur
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur lors de la réinitialisation du fichier")
	}

	// Créer le répertoire public si nécessaire
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.MkdirAll("public", os.ModePerm)
	}

	imageName := userID + "_" + file.Filename
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

	return ctx.JSON(http.StatusOK, map[string]string{
		"message":   "Image uploadée avec succès",
		"image_url": imagePath,
	})
}

func (c *UserController) GetUserEvents(ctx echo.Context) error {
	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	pagination := utils.PaginationFromContext(ctx)

	events, err := c.UserService.GetUserEvents(user.ID, pagination)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}

func (c *UserController) GetAssociationsEvents(ctx echo.Context) error {
	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	pagination := utils.PaginationFromContext(ctx)

	events, err := c.UserService.GetAssociationsEvents(user.ID, pagination)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}

func (c *UserController) ConfirmParticipation(ctx echo.Context) error {
	// Récupérer l'ID de la participation à confirmer depuis l'URL
	participationID := ctx.Param("id")

	// Récupérer l'utilisateur actuel depuis le contexte
	currentUser := ctx.Get("user")
	if currentUser == nil {
		return ctx.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Vérifiez que l'utilisateur est bien de type models.User
	user, ok := currentUser.(models.User)
	if !ok {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid user data"})
	}

	associationLeaderID := user.ID // ID de l'utilisateur actuel (leader de l'association)

	// Récupérer la participation avec les relations nécessaires
	var participation models.Participation
	if err := database.CurrentDatabase.Preload("Event.Association").
		Preload("Event.Category").
		Preload("User"). // Charger l'utilisateur lié à la participation
		First(&participation, "id = ?", participationID).Error; err != nil {
		ctx.Logger().Errorf("Participation not found: %v", err)
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Participation not found"})
	}

	// Vérifier que l'utilisateur est bien le leader de l'association de l'événement
	if participation.Event.Association.OwnerID != associationLeaderID {
		return ctx.JSON(http.StatusForbidden, map[string]string{"error": "You are not authorized to confirm this participation"})
	}

	// Vérifier si la participation est déjà confirmée ou déclinée
	if participation.Status != "pending" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Participation already handled"})
	}

	// Mettre à jour uniquement les champs spécifiques de la participation
	if err := database.CurrentDatabase.Model(&models.Participation{}).Where("id = ?", participationID).Updates(map[string]interface{}{
		"status":       "confirmed",
		"is_attending": true,
	}).Error; err != nil {
		ctx.Logger().Errorf("Failed to update participation: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update participation"})
	}

	// Incrémenter les points pour l'utilisateur lié à la participation
	if err := database.CurrentDatabase.Model(&models.User{}).Where("id = ?", participation.UserID).Update("points_open", gorm.Expr("points_open + ?", participation.Event.Category.Note)).Error; err != nil {
		ctx.Logger().Errorf("Failed to update user points: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user points"})
	}

	// Retourner une réponse réussie
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Participation confirmed successfully"})
}
