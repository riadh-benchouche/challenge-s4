package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) CreateMessage(message models.MessageCreate) (*models.Message, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(message); err != nil {
		return nil, err
	}

	newMessage := message.ToMessage()

	if err := database.CurrentDatabase.Create(newMessage).Error; err != nil {
		return nil, err
	}

	if err := database.CurrentDatabase.Preload("Sender").Preload("Association").
		First(newMessage, "id = ?", newMessage.ID).Error; err != nil {
		return nil, err
	}

	return newMessage, nil
}

func (s *MessageService) GetMessagesByAssociation(associationID string) ([]models.Message, error) {
	var messages []models.Message
	if err := database.CurrentDatabase.
		Preload("Sender").
		Preload("Association"). // Précharger les données de l'association
		Where("association_id = ?", associationID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) GetMessageByID(messageID string) (*models.Message, error) {
	var message models.Message
	if err := database.CurrentDatabase.Preload("Sender").Preload("Association").
		First(&message, "id = ?", messageID).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *MessageService) UpdateMessageContent(messageID string, updatedMessage models.MessageUpdate) (*models.Message, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(updatedMessage); err != nil {
		return nil, err
	}

	existingMessage := &models.Message{}
	if err := database.CurrentDatabase.First(existingMessage, "id = ?", messageID).Error; err != nil {
		return nil, err
	}

	if err := database.CurrentDatabase.Model(existingMessage).Updates(updatedMessage).Error; err != nil {
		return nil, err
	}

	if err := database.CurrentDatabase.Preload("Sender").Preload("Association").First(existingMessage, "id = ?", messageID).Error; err != nil {
		return nil, err
	}

	// Force une mise à jour du champ `updated_at` (bug GORM)
	if err := database.CurrentDatabase.Model(&models.Message{}).Where("id = ?", existingMessage.ID).Update("updated_at", time.Now()).Error; err != nil {
		return nil, err
	}

	return existingMessage, nil
}

// Supprime un message
func (s *MessageService) DeleteMessage(messageID string) error {
	if err := database.CurrentDatabase.Delete(&models.Message{}, "id = ?", messageID).Error; err != nil {
		return err
	}
	return nil
}

// Récupère les messages paginés d'une association
func (s *MessageService) GetMessagesByAssociationWithPagination(associationID string, pagination utils.Pagination) (*utils.Pagination, error) {
	var messages []models.Message
	query := database.CurrentDatabase.Where("association_id = ?", associationID).Order("created_at DESC")

	query.Scopes(utils.Paginate(messages, &pagination, query)).Preload("Sender").Find(&messages)
	pagination.Rows = messages

	return &pagination, nil
}
