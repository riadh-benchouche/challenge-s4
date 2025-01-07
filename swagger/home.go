package swagger

import (
	"backend/controllers"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupHomeSwagger(api *swag.API) {
	homeController := controllers.NewHomeController()

	// Endpoint: Get Statistics
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/statistics",
			endpoint.Handler(homeController.GetStatistics),
			endpoint.Summary("Retrieve user statistics"),
			endpoint.Description("Fetches statistics related to the logged-in user"),
			endpoint.Response(http.StatusOK, "Statistics data", endpoint.SchemaResponseOption(map[string]interface{}{
				"total_users":         "integer",
				"active_associations": "integer",
				"recent_events":       "integer",
			})),
			endpoint.Response(http.StatusUnauthorized, "User not authenticated"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Home"),
		),
	)

	// Endpoint: Get Top Associations
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/top-associations",
			endpoint.Handler(homeController.GetTopAssociations),
			endpoint.Summary("Retrieve top associations"),
			endpoint.Description("Fetches the top associations based on activity"),
			endpoint.Response(http.StatusOK, "List of top associations", endpoint.SchemaResponseOption([]map[string]interface{}{
				{
					"id":          "string",
					"name":        "string",
					"description": "string",
					"score":       "integer",
				},
			})),
			endpoint.Response(http.StatusUnauthorized, "User not authenticated"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Home"),
		),
	)

}
