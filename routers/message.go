package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupMessageRoutes(e *echo.Echo, controller *controllers.MessageController) {

	api := e.Group("/api")

	api.POST("/messages", controller.Create)
	api.GET("/messages", controller.GetAll)
	api.GET("/messages/:id", controller.GetByID)
	api.PUT("/messages/:id", controller.Update)
	api.DELETE("/messages/:id", controller.Delete)
}
