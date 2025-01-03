package swagger

import (
	"backend/controllers"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupAuthSwagger(api *swag.API) {
	authController := controllers.NewAuthController()

	// Login Endpoint
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/auth/login",
			endpoint.Handler(authController.Login),
			endpoint.Summary("User login"),
			endpoint.Description("Authenticates a user using their email and password"),
			endpoint.Body(map[string]string{
				"email":    "string (required)",
				"password": "string (required)",
			}, "Login credentials", true),
			endpoint.Response(http.StatusOK, "Login successful", endpoint.SchemaResponseOption(map[string]interface{}{
				"token": "string",
				"user":  "object (user details)",
			})),
			endpoint.Response(http.StatusUnauthorized, "Invalid credentials"),
			endpoint.Response(http.StatusUnprocessableEntity, "Validation error"),
			endpoint.Tags("Auth"),
		),
	)

	// Register Endpoint
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/auth/register",
			endpoint.Handler(authController.Register),
			endpoint.Summary("User registration"),
			endpoint.Description("Registers a new user with email, name, and password"),
			endpoint.Body(map[string]string{
				"email":         "string (required)",
				"password":      "string (required)",
				"name":          "string (required)",
				"plainPassword": "string (optional)",
			}, "Registration details", true),
			endpoint.Response(http.StatusCreated, "Registration successful", endpoint.SchemaResponseOption(map[string]interface{}{
				"user":  "object (user details)",
				"email": "string",
			})),
			endpoint.Response(http.StatusConflict, "Email already exists"),
			endpoint.Response(http.StatusUnprocessableEntity, "Validation error"),
			endpoint.Tags("Auth"),
		),
	)
}
