package controllers

import (
	"backend/enums"
	"backend/models"
	"backend/resources"
	"backend/services"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

type EventController struct {
	EventService       *services.EventService
	AssociationService *services.AssociationService
}

func NewEventController() *EventController {
	return &EventController{
		EventService:       services.NewEventService(),
		AssociationService: services.NewAssociationService(),
	}
}

func (c *EventController) CreateEvent(ctx echo.Context) error {
	var event models.Event
	err := json.NewDecoder(ctx.Request().Body).Decode(&event)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Mauvaise requête: Impossible de décoder le corps de la requête"})
	}

	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Non autorisé: Vous devez être connecté pour effectuer cette action"})
	}

	isUserInGroup, err := c.AssociationService.IsUserInAssociation(user.ID, event.AssociationID)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur serveur: Impossible de vérifier si l'utilisateur est membre de l'association"})
	}
	if !isUserInGroup {
		return ctx.JSON(http.StatusForbidden, map[string]string{"error": "Interdit: Vous n'êtes pas membre de cette association"})
	}

	newEvent, err := c.EventService.AddEvent(&event)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			validationErrors := utils.GetValidationErrors(validationErrs, event)
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"error":   "Validation error",
				"details": validationErrors,
			})
		}
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur serveur: Impossible de créer un évènement"})
	}

	return ctx.JSON(http.StatusCreated, newEvent)
}

func (c *EventController) GetEvents(ctx echo.Context) error {
	pagination := utils.PaginationFromContext(ctx)
	search := ctx.QueryParam("search")

	eventPagination, err := c.EventService.GetEvents(pagination, &search)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, eventPagination)
}

func (c *EventController) GetEventById(ctx echo.Context) error {
	id := ctx.Param("id")

	event, err := c.EventService.GetEventById(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Événement non trouvé")
	}

	return ctx.JSON(http.StatusOK, event)
}

func (c *EventController) UpdateEvent(ctx echo.Context) error {
	eventID := ctx.Param("id")

	existingEvent, err := c.EventService.GetEventById(eventID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Événement introuvable")
	}

	var updateData struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Date        *string `json:"date"`
		Location    *string `json:"location"`
	}

	if err := ctx.Bind(&updateData); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Données d'événement invalides")
	}

	if updateData.Name != nil {
		existingEvent.Name = *updateData.Name
	}
	if updateData.Description != nil {
		existingEvent.Description = *updateData.Description
	}
	if updateData.Date != nil {
		parsedDate, err := time.Parse(time.RFC3339, *updateData.Date)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, "Format de date invalide")
		}
		existingEvent.Date = parsedDate
	}
	if updateData.Location != nil {
		existingEvent.Location = *updateData.Location
	}

	if err := c.EventService.UpdateEvent(existingEvent); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}

	return ctx.JSON(http.StatusOK, existingEvent)
}

func (c *EventController) DeleteEvent(ctx echo.Context) error {
	id := ctx.Param("id")

	if _, err := ulid.Parse(id); err != nil {
		return ctx.JSON(http.StatusBadRequest, "ID invalide")
	}

	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.JSON(http.StatusUnauthorized, "Non autorisé")
	}

	event, err := c.EventService.GetEventById(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Événement non trouvé")
	}

	isUserInAssociation, err := c.AssociationService.IsUserInAssociation(user.ID, event.AssociationID)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, "Erreur serveur lors de la vérification de l'association")
	}

	association, err := c.AssociationService.GetAssociationById(event.AssociationID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur serveur lors de la récupération de l'association")
	}

	if !(user.Role == enums.AdminRole || (isUserInAssociation && association.OwnerID == user.ID)) {
		return ctx.JSON(http.StatusForbidden, "Interdit : vous n'avez pas les permissions nécessaires pour supprimer cet événement")
	}

	err = c.EventService.DeleteEvent(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Erreur serveur : impossible de supprimer l'événement")
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c *EventController) GetEventParticipations(ctx echo.Context) error {
	id := ctx.Param("id")

	if _, err := ulid.Parse(id); err != nil {
		return ctx.JSON(http.StatusBadRequest, "ID invalide")
	}

	pagination := utils.PaginationFromContext(ctx)

	statusParam := ctx.QueryParam("status")
	var status *string
	if statusParam != "" {
		status = &statusParam
	}

	participations, err := c.EventService.GetEventParticipations(id, pagination, status)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, "Erreur serveur : impossible de récupérer les participations")
	}

	participationResources := resources.NewParticipationResourceList(participations.Rows.([]models.Participation))

	pagination.Rows = participationResources

	return ctx.JSON(http.StatusOK, pagination)
}

func (c *EventController) GetUserEventParticipation(ctx echo.Context) error {
	user := ctx.Get("user").(models.User)

	eventID := ctx.Param("id")

	event, err := c.EventService.GetEventById(eventID)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if event == nil || event.ID == "" {
		return ctx.NoContent(http.StatusNotFound)
	}

	participation := c.EventService.GetUserEventParticipation(eventID, user.ID)
	if participation == nil {
		return ctx.JSON(http.StatusNotFound, "Participation introuvable")
	}

	return ctx.JSON(http.StatusOK, participation)
}

func (c *EventController) ChangeAttend(ctx echo.Context) error {
	user, ok := ctx.Get("user").(models.User)
	if !ok || user.ID == "" {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	eventID := ctx.Param("id")
	if eventID == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	event, err := c.EventService.GetEventById(eventID)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if event == nil || event.ID == "" {
		return ctx.NoContent(http.StatusNotFound)
	}

	body := struct {
		IsAttending bool `json:"is_attending"`
	}{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&body)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	isAttending := body.IsAttending

	participation, err := c.EventService.ChangeUserEventAttend(isAttending, eventID, user.ID)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, participation)
}
