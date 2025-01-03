package swagger

import (
	"backend/controllers"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupChatbotSwagger(api *swag.API) {
	chatbotController := controllers.NewChatbotController()

	// Endpoint: Chat with Chatbot
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/chatbot/message",
			endpoint.Handler(chatbotController.ChatHandler),
			endpoint.Summary("Interact with the chatbot"),
			endpoint.Description("Sends a message to the chatbot and retrieves its response."),
			endpoint.Body(map[string]string{
				"message": "string",
			}, "Chat request payload containing the user's message", true),
			endpoint.Response(http.StatusOK, "Chatbot response", endpoint.SchemaResponseOption(map[string]string{
				"response": "string",
			})),
			endpoint.Response(http.StatusBadRequest, "Invalid request payload"),
			endpoint.Response(http.StatusInternalServerError, "Error getting response from ChatGPT or searching associations"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Chatbot"),
		),
	)
}
