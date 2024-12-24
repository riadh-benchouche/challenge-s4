package routers

import (
	"backend/controllers"
	"backend/enums"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

type EventRouter struct{}

func (r EventRouter) SetupRoutes(e *echo.Echo) {
	eventController := controllers.NewEventController()
	api := e.Group("/events")

	api.GET("/:id/participation", eventController.GetUserEventParticipation, middlewares.AuthenticationMiddleware(enums.AdminRole, enums.AssociationLeaderRole, enums.UserRole))

	api.GET("/:id/participations", eventController.GetEventParticipations, middlewares.AuthenticationMiddleware(enums.AdminRole, enums.AssociationLeaderRole, enums.UserRole))
	api.POST("", eventController.CreateEvent, middlewares.AuthenticationMiddleware(enums.AssociationLeaderRole, enums.AdminRole))
	api.GET("", eventController.GetEvents, middlewares.AuthenticationMiddleware(enums.AdminRole))
	api.GET("/:id", eventController.GetEventById, middlewares.AuthenticationMiddleware(enums.AdminRole, enums.AssociationLeaderRole, enums.UserRole))
	api.PUT("/:id", eventController.UpdateEvent, middlewares.AuthenticationMiddleware(enums.AssociationLeaderRole, enums.AdminRole))
	api.DELETE("/:id", eventController.DeleteEvent, middlewares.AuthenticationMiddleware(enums.AdminRole, enums.AssociationLeaderRole))
	api.POST("/:id/user-event-participation", eventController.ChangeAttend, middlewares.AuthenticationMiddleware())
}
