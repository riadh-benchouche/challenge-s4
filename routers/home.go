package routers

import (
	"backend/controllers"
	"backend/middlewares"
	"github.com/labstack/echo/v4"
)

type HelloRouter struct{}

func (r *HelloRouter) SetupRoutes(e *echo.Echo) {
	homeController := controllers.NewHomeController()

	e.GET("/", homeController.Hello)
	e.GET("/admin/ping", homeController.HelloAdmin)
	e.GET("/statistics", homeController.GetStatistics, middlewares.AuthenticationMiddleware())
	e.GET("/top-associations", homeController.GetTopAssociations, middlewares.AuthenticationMiddleware())
}
