package middlewares

import (
	"backend/database"
	"backend/enums"
	"backend/models"
	"backend/services"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware(roles ...enums.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearer := c.Request().Header.Get("Authorization")

			if len(bearer) == 0 || strings.HasPrefix(bearer, "Bearer ") == false {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			token, err := services.NewAuthService().ValidateToken(bearer[7:]) // Supprime 'Bearer'
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			userID := claims["id"]
			var existingUser models.User
			database.CurrentDatabase.Preload("Groups").Where("id = ?", userID).First(&existingUser)
			if existingUser.ID == "" {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			if !existingUser.IsConfirmed || !existingUser.IsActive {
				return c.JSON(http.StatusUnauthorized, "Votre compte n'est pas actif ou confirmé")
			}

			if len(roles) > 0 {
				roleFound := false
				for _, role := range roles {
					if existingUser.Role == role {
						roleFound = true
						break
					}
				}

				if roleFound == false {
					return c.JSON(http.StatusUnauthorized, "unauthorized")
				}
			}

			c.Set("user", existingUser)

			return next(c)
		}
	}
}
