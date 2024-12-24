package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

type ChatbotRouter struct{}

func (r *ChatbotRouter) SetupRoutes(e *echo.Echo) {
	chatbotController := controllers.NewChatbotController()
	api := e.Group("/chatbot")
	api.POST("/message", chatbotController.ChatHandler)
}
