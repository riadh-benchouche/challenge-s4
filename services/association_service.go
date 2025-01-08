package services

import (
	"backend/database"
	"backend/enums"
	coreErrors "backend/errors"
	"backend/models"
	"backend/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AssociationService struct {
	db *gorm.DB
}

func NewAssociationService() *AssociationService {
	if database.CurrentDatabase == nil {
		return nil
	}
	return &AssociationService{
		db: database.CurrentDatabase,
	}
}

func (s *AssociationService) CreateAssociation(association models.Association) (*models.Association, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
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

func (s *AssociationService) GetAllAssociationsActiveAndNonActive(pagination utils.Pagination, filters ...AssociationFilter) (*utils.Pagination, error) {
	var associations []models.Association

	query := database.CurrentDatabase.
		Model(models.Association{})

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

func (s *AssociationService) GetAllAssociations(pagination utils.Pagination, filters ...AssociationFilter) (*utils.Pagination, error) {
	var associations []models.Association

	query := database.CurrentDatabase.
		Where("is_active = ?", true).
		Model(models.Association{})

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

func (s *AssociationService) GetNextEvent(groupID string) (*models.Event, error) {
	var event models.Event

	err := database.CurrentDatabase.
		Preload("Participations").
		Where("association_id = ?", groupID).
		Where("date >= ?", time.Now().Format(models.DateFormat)).
		Order("date").
		First(&event).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (s *AssociationService) GetAssociationEvents(groupID string, pagination utils.Pagination) (*utils.Pagination, error) {
	var events []models.Event

	query := database.CurrentDatabase.
		Preload("Participations").
		Where("association_id = ?", groupID).
		Where("date >= ?", time.Now().Format(models.DateFormat)).
		Order("date")

	query.Scopes(utils.Paginate(events, &pagination, query)).
		Find(&events)

	pagination.Rows = events

	return &pagination, nil
}

func (s *AssociationService) JoinAssociationByCode(userID string, code string) (*models.Association, error) {
	var association models.Association
	err := database.CurrentDatabase.Where("code = ?", code).First(&association).Error

	if err != nil {
		return nil, coreErrors.ErrAssociationNotFound
	}

	var membership models.Membership

	err = database.CurrentDatabase.Where("user_id = ? AND association_id = ?", userID, association.ID).First(&membership).Error
	if err == nil {
		return nil, coreErrors.ErrAlreadyJoined
	}

	NewMembership := models.Membership{
		UserID:        userID,
		AssociationID: association.ID,
		JoinedAt:      time.Now(),
		Status:        enums.Accepted,
	}

	err = database.CurrentDatabase.Create(&NewMembership).Error
	if err != nil {
		return nil, err
	}

	return &association, nil
}

func (s *AssociationService) UpdateAssociation(association *models.Association) error {
	var existingAssociation models.Association
	if err := database.CurrentDatabase.First(&existingAssociation, "id = ?", association.ID).Error; err != nil {
		return fmt.Errorf("Association not found: %w", err)
	}

	// Check for unique code conflict
	if association.Code != "" && association.Code != existingAssociation.Code {
		var otherAssociation models.Association
		if err := database.CurrentDatabase.First(&otherAssociation, "code = ?", association.Code).Error; err == nil {
			return fmt.Errorf("Code conflict: another association is already using this code")
		}
	}

	// Update fields
	updates := map[string]interface{}{
		"name":        association.Name,
		"description": association.Description,
		"is_active":   association.IsActive,
		"code":        association.Code,
		"image_url":   association.ImageURL,
	}

	if err := database.CurrentDatabase.Model(&existingAssociation).Updates(updates).Error; err != nil {
		return fmt.Errorf("Failed to update association: %w", err)
	}

	return nil
}
