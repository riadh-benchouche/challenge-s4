package swagger

import (
	"backend/controllers"
	"backend/models"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupMessageSwagger(api *swag.API) {
	messageController := controllers.NewMessageController()

	// Endpoint: Create Message
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/messages",
			endpoint.Handler(messageController.CreateMessage),
			endpoint.Summary("Create a new message"),
			endpoint.Description("Allows a user to create a new message in an association"),
			endpoint.Body(models.MessageCreate{}, "Message object to create", true),
			endpoint.Response(http.StatusCreated, "Successfully created message", endpoint.SchemaResponseOption(models.Message{})),
			endpoint.Response(http.StatusBadRequest, "Invalid input"),
			endpoint.Response(http.StatusUnauthorized, "User not authenticated"),
			endpoint.Response(http.StatusUnprocessableEntity, "Validation error"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Messages"),
		),
	)

	// Endpoint: Get Messages by Association
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/messages/association/{associationId}",
			endpoint.Handler(messageController.GetMessagesByAssociation),
			endpoint.Summary("Retrieve messages for an association"),
			endpoint.Description("Fetches all messages for a specific association"),
			endpoint.Path("associationId", "string", "ID of the association to fetch messages for", true),
			endpoint.Response(
				http.StatusOK,
				"List of messages",
				endpoint.SchemaResponseOption(map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/models.Message",
					},
				}),
			),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Messages"),
		),
	)

	// Endpoint: Update Message
	api.AddEndpoint(
		endpoint.New(
			http.MethodPut, "/messages/{id}",
			endpoint.Handler(messageController.UpdateMessage),
			endpoint.Summary("Update a message"),
			endpoint.Description("Allows a user to update the content of their message"),
			endpoint.Path("id", "string", "ID of the message to update", true),
			endpoint.Body(models.MessageUpdate{}, "Updated message content", true),
			endpoint.Response(http.StatusOK, "Successfully updated message", endpoint.SchemaResponseOption(models.Message{})),
			endpoint.Response(http.StatusBadRequest, "Invalid input or ID"),
			endpoint.Response(http.StatusUnauthorized, "User not authorized"),
			endpoint.Response(http.StatusNotFound, "Message not found"),
			endpoint.Response(http.StatusUnprocessableEntity, "Validation error"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Messages"),
		),
	)

	// Endpoint: Delete Message
	api.AddEndpoint(
		endpoint.New(
			http.MethodDelete, "/messages/{id}",
			endpoint.Handler(messageController.DeleteMessage),
			endpoint.Summary("Delete a message"),
			endpoint.Description("Allows a user to delete their message"),
			endpoint.Path("id", "string", "ID of the message to delete", true),
			endpoint.Response(http.StatusNoContent, "Successfully deleted message"),
			endpoint.Response(http.StatusBadRequest, "Invalid ID"),
			endpoint.Response(http.StatusUnauthorized, "User not authorized"),
			endpoint.Response(http.StatusNotFound, "Message not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Messages"),
		),
	)
}
