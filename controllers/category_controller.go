package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	service *services.CategoryService
}

func NewCategoryController(service *services.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// Create crée une nouvelle catégorie
func (c *CategoryController) Create(ctx echo.Context) error {
	category := new(models.Category)
	if err := ctx.Bind(category); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.service.Create(category); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, category)
}

// GetByID récupère une catégorie par son ID
func (c *CategoryController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	category, err := c.service.GetByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, category)
}

// GetAll récupère toutes les catégories
func (c *CategoryController) GetAll(ctx echo.Context) error {
	categories, err := c.service.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, categories)
}

// Update met à jour une catégorie
func (c *CategoryController) Update(ctx echo.Context) error {
	category := new(models.Category)
	if err := ctx.Bind(category); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	category.ID = ctx.Param("id")
	if err := c.service.Update(category); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, category)
}

// Delete supprime une catégorie
func (c *CategoryController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}
