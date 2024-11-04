package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupAssociationRoutes(e *echo.Echo, controller *controllers.AssociationController) {

	api := e.Group("/api")

	api.POST("/associations", controller.Create)
	api.GET("/associations", controller.GetAll)
	api.GET("/associations/:id", controller.GetByID)
	api.PUT("/associations/:id", controller.Update)
	api.DELETE("/associations/:id", controller.Delete)
}
