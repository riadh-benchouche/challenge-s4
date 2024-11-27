package routers

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct{}

func (r *AuthRouter) SetupRoutes(e *echo.Echo) {
	authController := controllers.NewAuthController()
	group := e.Group("/auth")

	// Routes existantes
	group.POST("/login", authController.Login)
	group.POST("/register", authController.Register)

	// Nouvelles routes pour la confirmation d'email
	group.GET("/confirm", authController.ConfirmEmail)
	group.POST("/resend-confirmation", authController.ResendConfirmation)
}

// Votre RegisterRequest reste inchangé
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
