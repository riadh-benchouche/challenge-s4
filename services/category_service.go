package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"errors"
	"strings"
)

type CategoryService struct {
}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	var existingCategory models.Category
	if err := database.CurrentDatabase.Where("name = ?", category.Name).First(&existingCategory).Error; err == nil {
		return errors.New("une catégorie avec ce nom existe déjà")
	}

	if err := database.CurrentDatabase.Create(category).Error; err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) GetCategories(pagination utils.Pagination, search *string) (*utils.Pagination, error) {
	var categories []models.Category
	query := database.CurrentDatabase

	if search != nil && *search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(*search)+"%")
	}

	err := query.Scopes(utils.Paginate(categories, &pagination, query)).
		Order("ID asc").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	pagination.Rows = categories
	return &pagination, nil
}

func (s *CategoryService) GetCategoryById(id string) (*models.Category, error) {
    var category models.Category
    if err := database.CurrentDatabase.First(&category, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &category, nil
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	var existingCategory models.Category
	if err := database.CurrentDatabase.First(&existingCategory, "id = ?", category.ID).Error; err != nil {
		return err
	}

	if err := database.CurrentDatabase.Model(&existingCategory).Updates(category).Error; err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) DeleteCategory(id string) error {
	if err := database.CurrentDatabase.Delete(&models.Category{}, "ID = ?", id).Error; err != nil {
		return err
	}
	return nil
}
