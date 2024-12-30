package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type MessageController struct {
	messageService   *services.MessageService
	webSocketService *services.WebSocketService
}

func NewMessageController() *MessageController {
	return &MessageController{
		messageService:   services.NewMessageService(),
		webSocketService: services.NewWebSocketService(),
	}
}

// Créer un nouveau message
func (c *MessageController) CreateMessage(ctx echo.Context) error {
	var jsonBody models.MessageCreate

	// Décodage de la requête JSON
	if err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	// Récupérer l'utilisateur authentifié
	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	// Ajouter l'utilisateur comme expéditeur
	jsonBody.SenderID = user.ID

	// Validation de l'entrée
	validate := validator.New()
	if err := validate.Struct(jsonBody); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, utils.GetValidationErrors(err.(validator.ValidationErrors), jsonBody))
	}

	// Créer le message
	newMessage, err := c.messageService.CreateMessage(jsonBody)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	// Diffuser le message via WebSocket
	if err := c.webSocketService.BroadcastToAssociation(newMessage.AssociationID, newMessage); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, newMessage)
}

// Mettre à jour le contenu d'un message
func (c *MessageController) UpdateMessage(ctx echo.Context) error {
	messageID := ctx.Param("id")

	var jsonBody models.MessageUpdate
	if err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	// Récupérer le message existant
	message, err := c.messageService.GetMessageByID(messageID)
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	// Vérifier que l'utilisateur est l'auteur du message
	if message.SenderID != user.ID {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	// Mettre à jour le contenu
	updatedMessage, err := c.messageService.UpdateMessageContent(messageID, jsonBody)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return ctx.JSON(http.StatusUnprocessableEntity, utils.GetValidationErrors(validationErrs, jsonBody))
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, updatedMessage)
}

// Supprimer un message
func (c *MessageController) DeleteMessage(ctx echo.Context) error {
	messageID := ctx.Param("id")

	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	// Récupérer le message
	message, err := c.messageService.GetMessageByID(messageID)
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	// Vérifier que l'utilisateur est l'auteur
	if message.SenderID != user.ID {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	// Supprimer le message
	if err := c.messageService.DeleteMessage(messageID); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Obtenir tous les messages d'une association
func (c *MessageController) GetMessagesByAssociation(ctx echo.Context) error {
	associationID := ctx.Param("associationId")

	// pagination := utils.PaginationFromContext(ctx)

	messages, err := c.messageService.GetMessagesByAssociation(associationID)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, messages)
}
