package routers

import (
	"backend/controllers"

	"backend/enums"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

type UserRouter struct{}

func (r *UserRouter) SetupRoutes(e *echo.Echo) {
	userController := controllers.NewUserController()

	group := e.Group("/users")
	group.POST("", userController.CreateUser, middlewares.AuthenticationMiddleware(enums.AdminRole))
	group.GET("", userController.GetUsers, middlewares.AuthenticationMiddleware(enums.AdminRole))
	group.PUT("/:id", userController.UpdateUser, middlewares.AuthenticationMiddleware())
	group.DELETE("/:id", userController.DeleteUser, middlewares.AuthenticationMiddleware(enums.AdminRole))
	group.GET("/:id", userController.FindByID, middlewares.AuthenticationMiddleware())
	group.GET("/:id/owner-associations", userController.GetOwnerAssociations, middlewares.AuthenticationMiddleware(enums.AssociationLeaderRole))
	group.GET("/:id/associations", userController.GetUserAssociations, middlewares.AuthenticationMiddleware())
	group.POST("/:id/associations/:association_id", userController.JoinAssociation, middlewares.AuthenticationMiddleware())
	group.POST("/:id/upload-image", userController.UploadProfileImage, middlewares.AuthenticationMiddleware())
	group.GET("/events", userController.GetUserEvents, middlewares.AuthenticationMiddleware())
	group.GET("/associations/events", userController.GetAssociationsEvents, middlewares.AuthenticationMiddleware())

	group.POST("/participations/:id/confirm", userController.ConfirmParticipation, middlewares.AuthenticationMiddleware(enums.AssociationLeaderRole))

}
