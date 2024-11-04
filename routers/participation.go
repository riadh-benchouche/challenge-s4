package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupParticipationRoutes(e *echo.Echo, controller *controllers.ParticipationController) {
	// Groupe de routes pour les participations
	api := e.Group("/api")

	// Routes CRUD
	api.POST("/participations", controller.Create)
	api.GET("/participations", controller.GetAll)
	api.GET("/participations/:id", controller.GetByID)
	api.PUT("/participations/:id", controller.Update)
	api.DELETE("/participations/:id", controller.Delete)
}
