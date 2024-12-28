package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketController struct {
	webSocketService *services.WebSocketService
}

func NewWebSocketController(webSocketService *services.WebSocketService) *WebSocketController {
	return &WebSocketController{
		webSocketService: webSocketService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permettre les connexions de toutes les origines (modifiable selon les besoins)
		return true
	},
}

func (controller *WebSocketController) OpenWebSocket(ctx echo.Context) error {
	// Récupération de l'utilisateur connecté depuis le contexte
	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "Unauthorized"})
	}

	// Établir la connexion WebSocket
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	// Gestion de la fermeture de la connexion
	defer controller.webSocketService.AcceptNewWebSocketConnection(ws, &user)(ws)

	// Gestion des messages WebSocket
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			ctx.Logger().Error(err)
			break
		}

		if err := controller.webSocketService.HandleWebSocketMessage(msg, &user, ws); err != nil {
			ctx.Logger().Warn(err)
			continue
		}
	}
	return nil
}
