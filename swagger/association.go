package swagger

import (
	"backend/controllers"
	"backend/models"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupAssociationSwagger(api *swag.API) {
	associationController := controllers.NewAssociationController()

	// Endpoint: Get All Associations
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/associations",
			endpoint.Handler(associationController.GetAllAssociations),
			endpoint.Summary("Retrieve all associations"),
			endpoint.Description("Fetches a list of all associations with optional filters and pagination"),
			endpoint.Query("filters", "string", "Filters in JSON format", false),
			endpoint.Query("page", "integer", "Page number for pagination", false),
			endpoint.Query("pageSize", "integer", "Number of items per page", false),
			endpoint.Response(http.StatusOK, "List of associations", endpoint.SchemaResponseOption([]models.Association{})),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Associations"),
		),
	)

	// Endpoint: Get Association By ID
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/associations/{associationId}",
			endpoint.Handler(associationController.GetAssociationById),
			endpoint.Summary("Retrieve an association by ID"),
			endpoint.Description("Fetches the details of a specific association using its ID"),
			endpoint.Path("associationId", "string", "ID of the association to retrieve", true),
			endpoint.Response(http.StatusOK, "Details of the association", endpoint.SchemaResponseOption(models.Association{})),
			endpoint.Response(http.StatusBadRequest, "Invalid ID format"),
			endpoint.Response(http.StatusNotFound, "Association not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Associations"),
		),
	)

	// Endpoint: Create Association
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/associations",
			endpoint.Handler(associationController.CreateAssociation),
			endpoint.Summary("Create a new association"),
			endpoint.Description("Allows a user to create a new association"),
			endpoint.Body(models.Association{}, "Association object to create", true),
			endpoint.Response(http.StatusCreated, "Successfully created association", endpoint.SchemaResponseOption(models.Association{})),
			endpoint.Response(http.StatusBadRequest, "Invalid input"),
			endpoint.Response(http.StatusUnauthorized, "User not authenticated"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Associations"),
		),
	)
}
