package services

import (
	"backend/models"

	"gorm.io/gorm"
)

type AssociationService struct {
	db *gorm.DB
}

func NewAssociationService(db *gorm.DB) *AssociationService {
	return &AssociationService{db: db}
}

func (s *AssociationService) Create(association *models.Association) error {
	return s.db.Create(association).Error
}

func (s *AssociationService) GetByID(id string) (*models.Association, error) {
	var association models.Association
	err := s.db.First(&association, "id = ?", id).Error
	return &association, err
}

func (s *AssociationService) GetAll() ([]models.Association, error) {
	var associations []models.Association
	err := s.db.Find(&associations).Error
	return associations, err
}

func (s *AssociationService) Update(association *models.Association) error {
	return s.db.Save(association).Error
}

func (s *AssociationService) Delete(id string) error {
	return s.db.Delete(&models.Association{}, "id = ?", id).Error
}
