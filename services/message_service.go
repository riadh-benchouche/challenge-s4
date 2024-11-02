package services

import (
	"backend/models"

	"gorm.io/gorm"
)

type MessageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) Create(message *models.Message) error {
	return s.db.Create(message).Error
}

func (s *MessageService) GetByID(id string) (*models.Message, error) {
	var message models.Message
	err := s.db.First(&message, "id = ?", id).Error
	return &message, err
}

func (s *MessageService) GetAll() ([]models.Message, error) {
	var messages []models.Message
	err := s.db.Find(&messages).Error
	return messages, err
}

func (s *MessageService) Update(message *models.Message) error {
	return s.db.Save(message).Error
}

func (s *MessageService) Delete(id string) error {
	return s.db.Delete(&models.Message{}, "id = ?", id).Error
}
