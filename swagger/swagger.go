package swagger

import (
	"github.com/zc2638/swag"
	"github.com/zc2638/swag/option"
)

// SetupSwagger configure tous les endpoints Swagger pour le projet
func SetupSwagger() *swag.API {
	api := swag.New(
		option.Title("API Documentation"),
		option.Description("Backend API documentation for all models"),
		option.Version("1.0.0"),
		option.SecurityScheme("bearer_auth",
			option.APIKeySecurity("Authorization", "header"),
		),
	)

	// Ajouter les endpoints pour chaque modèle
	SetupAssociationSwagger(api)
	SetupUserSwagger(api)
	SetupCategorySwagger(api)
	SetupEventSwagger(api)
	SetupMessageSwagger(api)
	SetupAuthSwagger(api)
	SetupChatbotSwagger(api)
	SetupHomeSwagger(api)
	// Ajouter d'autres endpoints ici pour d'autres modèles

	return api
}
