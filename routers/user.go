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
	group.POST("", userController.CreateUser)
	group.GET("", userController.GetUsers, middlewares.AuthenticationMiddleware(enums.AdminRole))
	group.PUT("/:id", userController.UpdateUser, middlewares.AuthenticationMiddleware())
	group.DELETE("/:id", userController.DeleteUser, middlewares.AuthenticationMiddleware(enums.AdminRole))
	group.GET("/:id", userController.FindByID, middlewares.AuthenticationMiddleware())
	// group.GET("/events", userController.GetUserEvents)
}
