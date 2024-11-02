package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupEventRoutes(e *echo.Echo, controller *controllers.EventController) {

	api := e.Group("/api")

	api.POST("/events", controller.Create)
	api.GET("/events", controller.GetAll)
	api.GET("/events/:id", controller.GetByID)
	api.PUT("/events/:id", controller.Update)
	api.DELETE("/events/:id", controller.Delete)
}
