package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ParticipationController struct {
	service *services.ParticipationService
}

func NewParticipationController(service *services.ParticipationService) *ParticipationController {
	return &ParticipationController{service: service}
}

// Create crée une nouvelle participation
func (c *ParticipationController) Create(ctx echo.Context) error {
	participation := new(models.Participation)
	if err := ctx.Bind(participation); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.service.Create(participation); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, participation)
}

// GetByID récupère une participation par son ID
func (c *ParticipationController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	participation, err := c.service.GetByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, participation)
}

// GetAll récupère toutes les participations
func (c *ParticipationController) GetAll(ctx echo.Context) error {
	participations, err := c.service.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, participations)
}

// Update met à jour une participation
func (c *ParticipationController) Update(ctx echo.Context) error {
	participation := new(models.Participation)
	if err := ctx.Bind(participation); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	participation.ID = ctx.Param("id")
	if err := c.service.Update(participation); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, participation)
}

// Delete supprime une participation
func (c *ParticipationController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}
