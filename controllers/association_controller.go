package controllers

import (
	"backend/models"   // Remplacez par votre import path
	"backend/services" // Remplacez par votre import path
	"net/http"

	"github.com/labstack/echo/v4"
)

type AssociationController struct {
	service *services.AssociationService
}

func NewAssociationController(service *services.AssociationService) *AssociationController {
	return &AssociationController{service: service}
}

// Create crée une nouvelle association
func (c *AssociationController) Create(ctx echo.Context) error {
	association := new(models.Association)
	if err := ctx.Bind(association); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.service.Create(association); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, association)
}

// GetByID récupère une association par son ID
func (c *AssociationController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	association, err := c.service.GetByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, association)
}

// GetAll récupère toutes les associations
func (c *AssociationController) GetAll(ctx echo.Context) error {
	associations, err := c.service.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, associations)
}

// Update met à jour une association
func (c *AssociationController) Update(ctx echo.Context) error {
	association := new(models.Association)
	if err := ctx.Bind(association); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	association.ID = ctx.Param("id")
	if err := c.service.Update(association); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, association)
}

// Delete supprime une association
func (c *AssociationController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}
