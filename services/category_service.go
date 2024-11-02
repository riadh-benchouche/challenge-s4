package services

import (
	"backend/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (s *CategoryService) Create(category *models.Category) error {
	return s.db.Create(category).Error
}

func (s *CategoryService) GetByID(id string) (*models.Category, error) {
	var category models.Category
	err := s.db.First(&category, "id = ?", id).Error
	return &category, err
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := s.db.Find(&categories).Error
	return categories, err
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.db.Save(category).Error
}

func (s *CategoryService) Delete(id string) error {
	return s.db.Delete(&models.Category{}, "id = ?", id).Error
}
