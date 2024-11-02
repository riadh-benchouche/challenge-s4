package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupMembershipRoutes(e *echo.Echo, controller *controllers.MembershipController) {
	api := e.Group("/api")

	api.POST("/memberships", controller.Create)
	api.GET("/memberships", controller.GetAll)
	api.GET("/memberships/:id", controller.GetByID)
	api.PUT("/memberships/:id", controller.Update)
	api.DELETE("/memberships/:id", controller.Delete)
}
