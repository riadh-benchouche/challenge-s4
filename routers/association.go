package routers

import (
	"backend/controllers"
	"backend/enums"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

type AssociationRouter struct{}

func (r *AssociationRouter) SetupRoutes(e *echo.Echo) {
	associationController := controllers.NewAssociationController()

	group := e.Group("/associations")

	// group.GET("", groupController.GetAllMyGroups, middlewares.AuthenticationMiddleware())
	group.GET("/all", associationController.GetAllAssociations, middlewares.AuthenticationMiddleware(enums.AdminRole))
	group.GET("/:associationId", associationController.GetAssociationById, middlewares.AuthenticationMiddleware(), middlewares.AssociationMembershipMiddleware)
	group.POST("", associationController.CreateAssociation, middlewares.AuthenticationMiddleware())
	// group.GET("/:groupId/next-event", groupController.GetNextEvent, middlewares.AuthenticationMiddleware(), middlewares.GroupMembershipMiddleware)
	// group.GET("/:groupId/events", groupController.GetGroupEvents, middlewares.AuthenticationMiddleware(), middlewares.GroupMembershipMiddleware)
}
