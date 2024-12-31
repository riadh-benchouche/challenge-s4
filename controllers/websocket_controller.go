package controllers

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketController struct {
	webSocketService *services.WebSocketService
}

// Initialisation du WebSocketController
func NewWebSocketController() *WebSocketController {
	return &WebSocketController{
		webSocketService: services.NewWebSocketService(),
	}
}

// Déclaration globale du WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permettre toutes les origines (modifiable selon les besoins)
		return true
	},
}

// Ouvrir une connexion WebSocket
func (controller *WebSocketController) OpenWebSocket(c echo.Context) error {
	// Récupération de l'utilisateur authentifié depuis le contexte Echo
	loggedUser, ok := c.Get("user").(models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Utilisateur non authentifié")
	}

	// Mise à niveau de la connexion HTTP vers une connexion WebSocket
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Impossible d'ouvrir la connexion WebSocket")
	}
	defer controller.webSocketService.AcceptNewWebSocketConnection(ws, &loggedUser)(ws)

	// Gestion des messages reçus via WebSocket
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error("Erreur lors de la lecture d'un message WebSocket: ", err)
			break
		}

		// Gestion du message via le WebSocketService
		if err := controller.webSocketService.HandleWebSocketMessage(msg, &loggedUser, ws); err != nil {
			c.Logger().Warn("Erreur lors du traitement du message WebSocket: ", err)
			continue
		}
	}

	return nil
}

// Diffuser un message à toutes les connexions
func (controller *WebSocketController) Broadcast(c echo.Context) error {
	// Exemple : données à diffuser passées via le corps de la requête
	var message models.Message
	if err := c.Bind(&message); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Données invalides")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur lors de la sérialisation du message")
	}

	if err := controller.webSocketService.Broadcast(data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur lors de la diffusion du message")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Message diffusé avec succès",
	})
}
