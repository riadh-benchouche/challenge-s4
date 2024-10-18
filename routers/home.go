package routers

import (
	"backend/controllers"
	"github.com/labstack/echo/v4"
)

type HelloRouter struct{}

func (r *HelloRouter) SetupRoutes(e *echo.Echo) {
	homeController := controllers.NewHomeController()

	e.GET("/", homeController.Hello)
	e.GET("/admin/ping", homeController.HelloAdmin)
}
