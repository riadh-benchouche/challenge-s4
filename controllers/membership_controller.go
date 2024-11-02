package controllers

import (
	"backend/models"   // Remplacez par votre import path
	"backend/services" // Remplacez par votre import path
	"net/http"

	"github.com/labstack/echo/v4"
)

type MembershipController struct {
	service *services.MembershipService
}

func NewMembershipController(service *services.MembershipService) *MembershipController {
	return &MembershipController{service: service}
}

// Create crée un nouveau membership
func (c *MembershipController) Create(ctx echo.Context) error {
	membership := new(models.Membership)
	if err := ctx.Bind(membership); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.service.Create(membership); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, membership)
}

// GetByID récupère un membership par son ID
func (c *MembershipController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	membership, err := c.service.GetByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, membership)
}

// GetAll récupère tous les memberships
func (c *MembershipController) GetAll(ctx echo.Context) error {
	memberships, err := c.service.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, memberships)
}

// Update met à jour un membership
func (c *MembershipController) Update(ctx echo.Context) error {
	membership := new(models.Membership)
	if err := ctx.Bind(membership); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	membership.ID = ctx.Param("id")
	if err := c.service.Update(membership); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, membership)
}

// Delete supprime un membership
func (c *MembershipController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}
