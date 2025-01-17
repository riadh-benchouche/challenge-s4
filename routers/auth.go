package routers

import (
	"backend/controllers"
	"backend/middlewares"
	"fmt"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct{}

func (r *AuthRouter) SetupRoutes(e *echo.Echo) {
	fmt.Println("🛣️ Setting up Auth routes...")

	authController := controllers.NewAuthController()
	group := e.Group("/auth")

	// Routes existantes
	group.POST("/login", authController.Login)
	group.POST("/register", authController.Register)
	group.POST("/refresh", authController.RefreshToken)
	fmt.Println("✅ Auth routes configured")

	// Routes de confirmation d'email avec le middleware de logging
	group.GET("/confirm", authController.ConfirmEmail, middlewares.EmailVerificationLoggingMiddleware)
	group.POST("/resend-confirmation", authController.ResendConfirmation)

	// Routes de mot de passe oublié et de réinitialisation
	group.POST("/forgot-password", authController.ForgotPassword)
	group.POST("/reset-password", authController.ResetPassword)
	group.GET("/reset-password", authController.ResetPasswordForm)

}
