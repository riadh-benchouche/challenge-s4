package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MessageController struct {
	service *services.MessageService
}

func NewMessageController(service *services.MessageService) *MessageController {
	return &MessageController{service: service}
}

// Create crée un nouveau message
func (c *MessageController) Create(ctx echo.Context) error {
	message := new(models.Message)
	if err := ctx.Bind(message); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.service.Create(message); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, message)
}

// GetByID récupère un message par son ID
func (c *MessageController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	message, err := c.service.GetByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, message)
}

// GetAll récupère tous les messages
func (c *MessageController) GetAll(ctx echo.Context) error {
	messages, err := c.service.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, messages)
}

// Update met à jour un message
func (c *MessageController) Update(ctx echo.Context) error {
	message := new(models.Message)
	if err := ctx.Bind(message); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	message.ID = ctx.Param("id")
	if err := c.service.Update(message); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, message)
}

// Delete supprime un message
func (c *MessageController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}
