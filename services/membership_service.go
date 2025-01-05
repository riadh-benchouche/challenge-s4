package services

import (
	"backend/models"

	"gorm.io/gorm"
)

type MembershipService struct {
	db *gorm.DB
}

func NewMembershipService(db *gorm.DB) *MembershipService {
	return &MembershipService{db: db}
}

func (s *MembershipService) Create(membership *models.Membership) error {
	return s.db.Create(membership).Error
}

func (s *MembershipService) GetByID(id string) (*models.Membership, error) {
	var membership models.Membership
	err := s.db.First(&membership, "id = ?", id).Error
	return &membership, err
}

func (s *MembershipService) GetAll() ([]models.Membership, error) {
	var memberships []models.Membership
	err := s.db.Find(&memberships).Error
	return memberships, err
}

func (s *MembershipService) Update(membership *models.Membership) error {
	return s.db.Save(membership).Error
}

func (s *MembershipService) Delete(id string) error {
	return s.db.Delete(&models.Membership{}, "id = ?", id).Error
}
