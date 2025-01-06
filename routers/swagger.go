package routers

import (
	"backend/swagger"

	"github.com/labstack/echo/v4"
	"github.com/zc2638/swag"
)

func SetupSwaggerRoutes(e *echo.Echo) {
	api := swagger.SetupSwagger()

	e.GET("/swagger/json", echo.WrapHandler(api.Handler()))
	e.GET("/swagger/ui*", echo.WrapHandler(swag.UIHandler("/swagger/ui", "/swagger/json", true)))
}
