package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupCategoryRoutes(e *echo.Echo, controller *controllers.CategoryController) {

	api := e.Group("/api")

	api.POST("/categories", controller.Create)
	api.GET("/categories", controller.GetAll)
	api.GET("/categories/:id", controller.GetByID)
	api.PUT("/categories/:id", controller.Update)
	api.DELETE("/categories/:id", controller.Delete)
}
