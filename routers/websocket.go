package routers

import (
	"backend/controllers"
	"backend/services"

	"github.com/labstack/echo/v4"
)

type WebSocketRouter struct {
	webSocketService *services.WebSocketService
}

func NewWebSocketRouter(webSocketService *services.WebSocketService) *WebSocketRouter {
	return &WebSocketRouter{
		webSocketService: webSocketService,
	}
}

func (r *WebSocketRouter) SetupRoutes(e *echo.Echo) {
	// Créez le contrôleur en passant WebSocketService
	webSocketController := controllers.NewWebSocketController(r.webSocketService)
	e.GET("/ws", webSocketController.OpenWebSocket)
}
