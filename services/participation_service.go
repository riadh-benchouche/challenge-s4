package services

import (
	"backend/models"
	"fmt"

	"gorm.io/gorm"
)

type ParticipationService struct {
	db *gorm.DB
}

func NewParticipationService(db *gorm.DB) *ParticipationService {
	return &ParticipationService{db: db}
}

func (s *ParticipationService) Create(participation *models.Participation) error {

	var event models.Event
	if err := s.db.First(&event, "id = ?", participation.EventID).Error; err != nil {
		return fmt.Errorf("event not found")
	}

	return s.db.Create(participation).Error
}

func (s *ParticipationService) GetByID(id string) (*models.Participation, error) {
	var participation models.Participation
	err := s.db.First(&participation, "id = ?", id).Error
	return &participation, err
}

func (s *ParticipationService) GetAll() ([]models.Participation, error) {
	var participations []models.Participation
	err := s.db.Find(&participations).Error
	return participations, err
}

func (s *ParticipationService) Update(participation *models.Participation) error {
	return s.db.Save(participation).Error
}

func (s *ParticipationService) Delete(id string) error {
	return s.db.Delete(&models.Participation{}, "id = ?", id).Error
}
