package routers

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

type CategoryRouter struct{}

func (r *CategoryRouter) SetupRoutes(e *echo.Echo) {
	categoryController := controllers.NewCategoryController()

	categoryGroup := e.Group("/categories")

	categoryGroup.POST("", categoryController.CreateCategory, middlewares.AuthenticationMiddleware())
	categoryGroup.GET("", categoryController.GetCategories, middlewares.AuthenticationMiddleware())
	// categoryGroup.GET("/:id", categoryController.GetCategoryByID)
	// categoryGroup.PUT("/:id", categoryController.UpdateCategory, middlewares.AuthenticationMiddleware())
	// categoryGroup.DELETE("/:id", categoryController.DeleteCategory, middlewares.AuthenticationMiddleware())
}
