package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"strings"

	"github.com/go-playground/validator/v10"
)

type EventService struct {
}

func NewEventService() *EventService {
	return &EventService{}
}

func (s *EventService) AddEvent(event *models.Event) (*models.Event, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(event)
	if err != nil {
		return nil, err
	}

	create := database.CurrentDatabase.Create(event)
	if create.Error != nil {
		return nil, create.Error
	}

	return event, nil
}

func (s *EventService) GetEvents(pagination utils.Pagination, search *string) (*utils.Pagination, error) {
	var events []models.Event
	query := database.CurrentDatabase

	if search != nil && *search != "" {
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			"%"+strings.ToLower(*search)+"%",
			"%"+strings.ToLower(*search)+"%",
		)
	}

	err := query.Scopes(utils.Paginate(events, &pagination, query)).
		Order("date asc").
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	pagination.Rows = events
	return &pagination, nil
}

func (s *EventService) GetEventById(id string) (*models.Event, error) {
	var event models.Event
	if err := database.CurrentDatabase.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *EventService) UpdateEvent(event *models.Event) error {
	var existingEvent models.Event
	if err := database.CurrentDatabase.First(&existingEvent, "id = ?", event.ID).Error; err != nil {
		return err
	}

	if err := database.CurrentDatabase.Model(&existingEvent).Updates(map[string]interface{}{
		"name":        event.Name,
		"description": event.Description,
		"date":        event.Date,
		"location":    event.Location,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (s *EventService) DeleteEvent(id string) error {
	if err := database.CurrentDatabase.Delete(&models.Event{}, "ID = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (s *EventService) GetEventParticipations(eventID string, pagination utils.Pagination, status *string) (*utils.Pagination, error) {
	var participations []models.Participation

	query := database.CurrentDatabase.
		Where("event_id = ?", eventID).
		Order("updated_at DESC")

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	query.Preload("User").
		Scopes(utils.Paginate(participations, &pagination, query)).
		Find(&participations)

	pagination.Rows = participations
	return &pagination, nil
}

func (s *EventService) GetUserEventParticipation(eventID string, userID string, preloads ...string) *models.Participation {
	var participation models.Participation
	query := database.CurrentDatabase.Where("event_id = ? AND user_id = ?", eventID, userID)

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	query.First(&participation)

	if participation.UserID == "" {
		return nil
	}

	return &participation
}

func (s *EventService) ChangeUserEventAttend(isAttending bool, eventID string, userID string) (*models.Participation, error) {
	participation := s.GetUserEventParticipation(eventID, userID)
	if participation == nil {
		participation = &models.Participation{
			EventID:     eventID,
			UserID:      userID,
			IsAttending: isAttending,
		}

		err := database.CurrentDatabase.Create(&participation).Error
		if err != nil {
			return nil, err
		}
		return participation, nil
	} else {
		participation.IsAttending = isAttending
		err := database.CurrentDatabase.Save(&participation).Error
		if err != nil {
			return nil, err
		}
		return participation, nil
	}
}
