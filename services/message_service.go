package services

import (
	"backend/database"
	"backend/models"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) Create(message *models.Message) error {
	return database.CurrentDatabase.Create(message).Error
}

func (s *MessageService) GetByID(id string) (*models.Message, error) {
	var message models.Message
	err := database.CurrentDatabase.First(&message, "id = ?", id).Error
	return &message, err
}

func (s *MessageService) GetAll() ([]models.Message, error) {
	var messages []models.Message
	err := database.CurrentDatabase.Find(&messages).Error
	return messages, err
}

func (s *MessageService) Update(message *models.Message) error {
	return database.CurrentDatabase.Save(message).Error
}

func (s *MessageService) Delete(id string) error {
	return database.CurrentDatabase.Delete(&models.Message{}, "id = ?", id).Error
}
