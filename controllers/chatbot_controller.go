package controllers

import (
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatbotController struct {
	ChatService *services.ChatService
}

func NewChatbotController() *ChatbotController {
	return &ChatbotController{
		ChatService: services.NewChatService(),
	}
}

func (cc *ChatbotController) ChatHandler(c echo.Context) error {
	var chatReq services.ChatRequest
	if err := c.Bind(&chatReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Utilisation du service pour obtenir la réponse de ChatGPT
	response, err := cc.ChatService.GetChatGPTResponse(chatReq.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting response from ChatGPT"})
	}

	return c.JSON(http.StatusOK, map[string]string{"response": response})
}
