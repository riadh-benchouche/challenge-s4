package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

type CategoryController struct {
	CategoryService *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		CategoryService: services.NewCategoryService(),
	}
}

func (c *CategoryController) CreateCategory(ctx echo.Context) error {
	var category models.Category
	if err := ctx.Bind(&category); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Données de catégorie invalides")
	}

	if category.Name == "" {
		return ctx.JSON(http.StatusBadRequest, "Le nom de la catégorie est requis")
	}

	category.ID = utils.GenerateULID()
	category.Note = 0
	if _, err := ulid.Parse(category.ID); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Format ULID invalide")
	}

	if err := c.CategoryService.CreateCategory(&category); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}

	return ctx.JSON(http.StatusCreated, category)
}

func (c *CategoryController) GetCategories(ctx echo.Context) error {
	pagination := utils.PaginationFromContext(ctx)
	search := ctx.QueryParam("search")

	categoryPagination, err := c.CategoryService.GetCategories(pagination, &search)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, categoryPagination)
}
