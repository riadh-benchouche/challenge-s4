package controllers

import (
	"backend/services"
	"fmt"
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
	fmt.Println("ChatHandler reached")

	var chatReq services.ChatRequest
	if err := c.Bind(&chatReq); err != nil {
		fmt.Println("Error binding request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Recherche dans la base de données
	assocService := services.NewAssociationService()
	dbData, err := assocService.SearchAssociations(chatReq.Message)
	if err != nil {
		fmt.Println("Error searching associations:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error searching associations"})
	}

	// Obtenez une réponse de ChatGPT enrichie
	response, err := cc.ChatService.GetChatGPTResponse(chatReq.Message, dbData)
	if err != nil {
		fmt.Println("Error getting ChatGPT response:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting response from ChatGPT"})
	}

	return c.JSON(http.StatusOK, map[string]string{"response": response})
}
