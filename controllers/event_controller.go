package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EventController struct {
	service *services.EventService
}

func NewEventController(service *services.EventService) *EventController {
	return &EventController{service: service}
}

// Create crée un nouvel événement
func (c *EventController) Create(ctx echo.Context) error {
	event := new(models.Event)
	if err := ctx.Bind(event); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.service.Create(event); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, event)
}

// GetByID récupère un événement par son ID
func (c *EventController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	event, err := c.service.GetByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, event)
}

// GetAll récupère tous les événements
func (c *EventController) GetAll(ctx echo.Context) error {
	events, err := c.service.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, events)
}

// Update met à jour un événement
func (c *EventController) Update(ctx echo.Context) error {
	event := new(models.Event)
	if err := ctx.Bind(event); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	event.ID = ctx.Param("id")
	if err := c.service.Update(event); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, event)
}

// Delete supprime un événement
func (c *EventController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}
