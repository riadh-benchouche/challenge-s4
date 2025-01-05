package routers

import "github.com/labstack/echo/v4"

type Router interface {
	SetupRoutes(e *echo.Echo)
}

func LoadRoutes(e *echo.Echo, routers ...Router) {
	for _, r := range routers {
		r.SetupRoutes(e)
	}
}
