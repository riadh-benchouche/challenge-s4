package routers

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

type MessageRouter struct{}

func (r *MessageRouter) SetupRoutes(e *echo.Echo) {
	messageController := controllers.NewMessageController()

	messageGroup := e.Group("/messages", middlewares.AuthenticationMiddleware())

	messageGroup.POST("", messageController.CreateMessage)
	messageGroup.GET("/association/:associationId", messageController.GetMessagesByAssociation)
	messageGroup.PUT("/:id", messageController.UpdateMessage)
	messageGroup.DELETE("/:id", messageController.DeleteMessage)
}
