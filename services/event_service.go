package services

import (
	"backend/models"

	"gorm.io/gorm"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

func (s *EventService) Create(event *models.Event) error {
	return s.db.Create(event).Error
}

func (s *EventService) GetByID(id string) (*models.Event, error) {
	var event models.Event
	err := s.db.First(&event, "id = ?", id).Error
	return &event, err
}

func (s *EventService) GetAll() ([]models.Event, error) {
	var events []models.Event
	err := s.db.Find(&events).Error
	return events, err
}

func (s *EventService) Update(event *models.Event) error {
	return s.db.Save(event).Error
}

func (s *EventService) Delete(id string) error {
	return s.db.Delete(&models.Event{}, "id = ?", id).Error
}
