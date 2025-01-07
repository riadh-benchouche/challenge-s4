package swagger

import (
	"backend/controllers"
	"backend/models"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupEventSwagger(api *swag.API) {
	eventController := controllers.NewEventController()

	// Endpoint: Create Event
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/events",
			endpoint.Handler(eventController.CreateEvent),
			endpoint.Summary("Create a new event"),
			endpoint.Description("Allows an authorized user to create a new event"),
			endpoint.Body(models.Event{}, "Event object to create", true),
			endpoint.Response(http.StatusCreated, "Successfully created event", endpoint.SchemaResponseOption(models.Event{})),
			endpoint.Response(http.StatusBadRequest, "Invalid event data"),
			endpoint.Response(http.StatusUnauthorized, "User not authenticated"),
			endpoint.Response(http.StatusForbidden, "User not authorized"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Tags("Events"),
		),
	)

	// Endpoint: Get All Events
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/events",
			endpoint.Handler(eventController.GetEvents),
			endpoint.Summary("Retrieve all events"),
			endpoint.Description("Fetches a paginated list of all events"),
			endpoint.Query("search", "string", "Optional search query", false),
			endpoint.Query("page", "integer", "Page number for pagination", false),
			endpoint.Query("pageSize", "integer", "Number of items per page", false),
			endpoint.Response(
				http.StatusOK,
				"List of events",
				endpoint.SchemaResponseOption(map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/models.Event",
					},
				}),
			),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Tags("Events"),
		),
	)

	// Endpoint: Get Event By ID
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/events/{id}",
			endpoint.Handler(eventController.GetEventById),
			endpoint.Summary("Retrieve an event by ID"),
			endpoint.Description("Fetches the details of a specific event using its ID"),
			endpoint.Path("id", "string", "ID of the event to retrieve", true),
			endpoint.Response(http.StatusOK, "Details of the event", endpoint.SchemaResponseOption(models.Event{})),
			endpoint.Response(http.StatusNotFound, "Event not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Tags("Events"),
		),
	)

	// Endpoint: Update Event
	api.AddEndpoint(
		endpoint.New(
			http.MethodPut, "/events/{id}",
			endpoint.Handler(eventController.UpdateEvent),
			endpoint.Summary("Update an event"),
			endpoint.Description("Allows an authorized user to update an existing event"),
			endpoint.Path("id", "string", "ID of the event to update", true),
			endpoint.Body(models.Event{}, "Updated event data", true),
			endpoint.Response(http.StatusOK, "Successfully updated event", endpoint.SchemaResponseOption(models.Event{})),
			endpoint.Response(http.StatusBadRequest, "Invalid event data"),
			endpoint.Response(http.StatusNotFound, "Event not found"),
			endpoint.Response(http.StatusConflict, "Event update conflict"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Tags("Events"),
		),
	)

	// Endpoint: Delete Event
	api.AddEndpoint(
		endpoint.New(
			http.MethodDelete, "/events/{id}",
			endpoint.Handler(eventController.DeleteEvent),
			endpoint.Summary("Delete an event"),
			endpoint.Description("Allows an authorized user to delete an event by ID"),
			endpoint.Path("id", "string", "ID of the event to delete", true),
			endpoint.Response(http.StatusNoContent, "Successfully deleted event"),
			endpoint.Response(http.StatusBadRequest, "Invalid ID"),
			endpoint.Response(http.StatusNotFound, "Event not found"),
			endpoint.Response(http.StatusForbidden, "User not authorized"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Tags("Events"),
		),
	)

	// Endpoint: Get Event Participations
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/events/{id}/participations",
			endpoint.Handler(eventController.GetEventParticipations),
			endpoint.Summary("Get participations for an event"),
			endpoint.Description("Fetches a list of users participating in the specified event"),
			endpoint.Path("id", "string", "ID of the event", true),
			endpoint.Response(
				http.StatusOK,
				"List of participations",
				endpoint.SchemaResponseOption(map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/models.Participation",
					},
				}),
			),
			endpoint.Response(http.StatusBadRequest, "Invalid ID"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Tags("Events"),
		),
	)

	// Endpoint: Change Attendance
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/events/{id}/user-event-participation",
			endpoint.Handler(eventController.ChangeAttend),
			endpoint.Summary("Change user attendance for an event"),
			endpoint.Description("Allows a user to mark their attendance for an event"),
			endpoint.Path("id", "string", "ID of the event", true),
			endpoint.Body(map[string]bool{"is_attending": true}, "Attendance status", true),
			endpoint.Response(http.StatusOK, "Attendance updated successfully", endpoint.SchemaResponseOption(models.Participation{})),
			endpoint.Response(http.StatusBadRequest, "Invalid request"),
			endpoint.Response(http.StatusUnauthorized, "User not authorized"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Tags("Events"),
		),
	)
}
