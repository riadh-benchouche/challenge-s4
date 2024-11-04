package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
	// "challenge-s4/middlewares"
	// "backend/models"
)

type UserRouter struct{}

func (r *UserRouter) SetupRoutes(e *echo.Echo) {
	userController := controllers.NewUserController()

	group := e.Group("/users")
	group.POST("", userController.CreateUser)
	group.GET("", userController.GetUsers)
	// group.PUT("/:id", userController.UpdateUser)
	// group.DELETE("/:id", userController.DeleteUser)
	// group.GET("/:id", userController.FindByID)
	// group.GET("/events", userController.GetUserEvents)
}
