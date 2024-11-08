package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AssociationService struct {
	db *gorm.DB
}

func NewAssociationService() *AssociationService {
	return &AssociationService{}
}

func (s *AssociationService) CreateAssociation(association models.Association) (*models.Association, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(association); err != nil {
		return nil, err
	}

	newAssociation := association.ToAssociation()
	if err := database.CurrentDatabase.Create(newAssociation).Error; err != nil {
		return nil, err
	}

	return newAssociation, nil
}

func (s *AssociationService) IsUserInAssociation(userId string, associationId string) (bool, error) {
	var association models.Association

	err := database.CurrentDatabase.Joins(
		"JOIN memberships ON memberships.association_id = associations.id AND memberships.user_id = ?", userId,
	).Where("associations.id = ?", associationId).First(&association).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return association.ID != "", nil
}

func (s *AssociationService) GetAssociationById(id string) (*models.Association, error) {
	var association models.Association
	if err := database.CurrentDatabase.Preload("Members").Preload("Owner").First(&association, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &association, nil
}

type AssociationFilter struct {
	database.Filter
	Column string `json:"column" validate:"required,oneof=name code"`
}

func (s *AssociationService) GetAllAssociations(pagination utils.Pagination, filters ...AssociationFilter) (*utils.Pagination, error) {
	var associations []models.Association

	query := database.CurrentDatabase.Model(models.Association{})

	if len(filters) > 0 {
		for _, filter := range filters {
			query = query.Where(filter.Column+" ILIKE ?", "%"+fmt.Sprintf("%v", filter.Value)+"%")
		}
	}

	err := query.Scopes(utils.Paginate(associations, &pagination, query)).
		Find(&associations).Error

	if err != nil {
		return nil, err
	}

	pagination.Rows = associations

	return &pagination, nil
}

func (s *AssociationService) Update(association *models.Association) error {
	return s.db.Save(association).Error
}

func (s *AssociationService) Delete(id string) error {
	return s.db.Delete(&models.Association{}, "id = ?", id).Error
}
