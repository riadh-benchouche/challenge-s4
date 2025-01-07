package swagger

import (
	"backend/controllers"
	"backend/models"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupUserSwagger(api *swag.API) {
	userController := controllers.NewUserController()

	// Endpoint: Create User
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/users",
			endpoint.Handler(userController.CreateUser),
			endpoint.Summary("Create a new user"),
			endpoint.Description("Allows an admin to create a new user"),
			endpoint.Body(models.User{}, "User object to create", true),
			endpoint.Response(http.StatusCreated, "Successfully created user", endpoint.SchemaResponseOption(models.User{})),
			endpoint.Response(http.StatusBadRequest, "Invalid input"),
			endpoint.Response(http.StatusConflict, "User already exists"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)

	// Endpoint: Get Users
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/users",
			endpoint.Handler(userController.GetUsers),
			endpoint.Summary("Retrieve all users"),
			endpoint.Description("Fetches a paginated list of users with optional filters"),
			endpoint.Query("filters", "string", "Filters in JSON format", false),
			endpoint.Query("page", "integer", "Page number for pagination", false),
			endpoint.Query("pageSize", "integer", "Number of items per page", false),
			endpoint.Response(http.StatusOK, "List of users", endpoint.SchemaResponseOption([]models.User{})),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)

	// Endpoint: Find User By ID
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/users/{id}",
			endpoint.Handler(userController.FindByID),
			endpoint.Summary("Retrieve a user by ID"),
			endpoint.Description("Fetches details of a specific user by ID"),
			endpoint.Path("id", "string", "ID of the user to retrieve", true),
			endpoint.Response(http.StatusOK, "Details of the user", endpoint.SchemaResponseOption(models.User{})),
			endpoint.Response(http.StatusBadRequest, "Invalid ID format"),
			endpoint.Response(http.StatusForbidden, "Unauthorized access"),
			endpoint.Response(http.StatusNotFound, "User not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)

	// Endpoint: Update User
	api.AddEndpoint(
		endpoint.New(
			http.MethodPut, "/users/{id}",
			endpoint.Handler(userController.UpdateUser),
			endpoint.Summary("Update a user"),
			endpoint.Description("Allows a user to update their profile or an admin to update user information"),
			endpoint.Path("id", "string", "ID of the user to update", true),
			endpoint.Body(models.User{}, "Updated user object", true),
			endpoint.Response(http.StatusOK, "Successfully updated user", endpoint.SchemaResponseOption(models.User{})),
			endpoint.Response(http.StatusBadRequest, "Invalid input or ID format"),
			endpoint.Response(http.StatusForbidden, "Unauthorized access"),
			endpoint.Response(http.StatusNotFound, "User not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)

	// Endpoint: Delete User
	api.AddEndpoint(
		endpoint.New(
			http.MethodDelete, "/users/{id}",
			endpoint.Handler(userController.DeleteUser),
			endpoint.Summary("Delete a user"),
			endpoint.Description("Allows an admin to delete a user by ID"),
			endpoint.Path("id", "string", "ID of the user to delete", true),
			endpoint.Response(http.StatusNoContent, "Successfully deleted user"),
			endpoint.Response(http.StatusBadRequest, "Invalid ID format"),
			endpoint.Response(http.StatusForbidden, "Cannot delete your own account"),
			endpoint.Response(http.StatusNotFound, "User not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)

	// Endpoint: Get User Events
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/users/events",
			endpoint.Handler(userController.GetUserEvents),
			endpoint.Summary("Retrieve user-specific events"),
			endpoint.Description("Fetches all events associated with the logged-in user"),
			endpoint.Query("page", "integer", "Page number for pagination", false),
			endpoint.Query("pageSize", "integer", "Number of items per page", false),
			endpoint.Response(http.StatusOK, "List of user events", endpoint.SchemaResponseOption([]models.Event{})),
			endpoint.Response(http.StatusUnauthorized, "Unauthorized access"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)

	// Endpoint: Get Associations Events
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/users/associations/events",
			endpoint.Handler(userController.GetAssociationsEvents),
			endpoint.Summary("Retrieve events for associations"),
			endpoint.Description("Fetches all events associated with the user's associations"),
			endpoint.Query("page", "integer", "Page number for pagination", false),
			endpoint.Query("pageSize", "integer", "Number of items per page", false),
			endpoint.Response(http.StatusOK, "List of association events", endpoint.SchemaResponseOption([]models.Event{})),
			endpoint.Response(http.StatusUnauthorized, "Unauthorized access"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)

	// Endpoint: Join Association
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/users/{id}/associations/{association_id}",
			endpoint.Handler(userController.JoinAssociation),
			endpoint.Summary("Join an association"),
			endpoint.Description("Allows a user to join an association using a code"),
			endpoint.Path("id", "string", "ID of the user", true),
			endpoint.Path("association_id", "string", "ID of the association", true),
			endpoint.Query("code", "string", "Code to join the association", true),
			endpoint.Response(http.StatusCreated, "User successfully joined the association", endpoint.SchemaResponseOption(map[string]string{
				"message": "User successfully joined the association",
			})),
			endpoint.Response(http.StatusBadRequest, "Invalid ULID format or missing code"),
			endpoint.Response(http.StatusConflict, "User already joined"),
			endpoint.Response(http.StatusUnauthorized, "Invalid association code"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Users"),
		),
	)
}
