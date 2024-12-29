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

	group.GET("", associationController.GetAllAssociations, middlewares.AuthenticationMiddleware())
	group.GET("/:associationId", associationController.GetAssociationById, middlewares.AuthenticationMiddleware(), middlewares.AssociationMembershipMiddleware)
	group.POST("", associationController.CreateAssociation, middlewares.AuthenticationMiddleware(enums.AssociationLeaderRole))
	group.POST("/:id/upload-image", associationController.UploadProfileImage, middlewares.AuthenticationMiddleware())
	group.GET("/:associationId/next-event", associationController.GetNextEvent, middlewares.AuthenticationMiddleware(), middlewares.AssociationMembershipMiddleware)
	group.GET("/:associationId/events", associationController.GetAssociationEvents, middlewares.AuthenticationMiddleware(), middlewares.AssociationMembershipMiddleware)
}
