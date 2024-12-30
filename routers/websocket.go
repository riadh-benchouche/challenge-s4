package routers

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

type WebSocketRouter struct{}

func (r *WebSocketRouter) SetupRoutes(e *echo.Echo) {
	wsController := controllers.NewWebSocketController()

	// Route pour WebSocket
	e.GET("/ws", wsController.OpenWebSocket, middlewares.AuthenticationMiddleware())
}
