package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct{}

func (r *AuthRouter) SetupRoutes(e *echo.Echo) {
	authController := controllers.NewAuthController()

	group := e.Group("/auth")
	group.POST("/login", authController.Login)
	group.POST("/register", authController.Register)
}
