package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func EmailVerificationLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Log uniquement pour les routes de vérification d'email
		if c.Path() == "/verify-email" {
			token := c.QueryParam("token")
			fmt.Printf(" Email verification request received\n")
			fmt.Printf(" Token: %s\n", token)
			fmt.Printf(" Request Method: %s\n", c.Request().Method)
			fmt.Printf(" Full URL: %s\n", c.Request().URL.String())
		}

		return next(c)
	}
}
