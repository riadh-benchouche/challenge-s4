package middlewares

import (
	"backend/enums"
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

// AssociationMembershipMiddleware vérifie que l'utilisateur appartient à une association
func AssociationMembershipMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		associationID := c.Param("associationId")

		_, err := ulid.Parse(associationID)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user, ok := c.Get("user").(models.User)
		if !ok || user.ID == "" {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}

		isUserInAssociation, err := services.NewAssociationService().IsUserInAssociation(user.ID, associationID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}

		// Si l'utilisateur n'est pas membre et n'est pas administrateur, refuser l'accès
		if !isUserInAssociation && !enums.IsAdmin(user.Role) {
			return c.JSON(http.StatusUnauthorized, "user does not belong to this association")
		}

		return next(c)
	}
}
