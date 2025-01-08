package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ParticipationController struct {
	service *services.ParticipationService
}

func NewParticipationController(service *services.ParticipationService) *ParticipationController {
	return &ParticipationController{service: service}
}

// Dans participation_controller.go
func (c *ParticipationController) CreateParticipation(ctx echo.Context) error {
	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var participation models.Participation
	if err := ctx.Bind(&participation); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	participation.UserID = user.ID

	err := c.service.Create(&participation)
	if err != nil {
		if err.Error() == "event not found" {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "event not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, participation)
}

func (c *ParticipationController) GetParticipations(ctx echo.Context) error {
	_, ok := ctx.Get("user").(models.User)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	participations, err := c.service.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, participations)
}

func (c *ParticipationController) UpdateParticipation(ctx echo.Context) error {
	_, ok := ctx.Get("user").(models.User)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	var participation models.Participation
	if err := ctx.Bind(&participation); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := c.service.Update(&participation)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, participation)
}

func (c *ParticipationController) DeleteParticipation(ctx echo.Context) error {
	_, ok := ctx.Get("user").(models.User)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	err := c.service.Delete(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "participation deleted successfully"})
}
