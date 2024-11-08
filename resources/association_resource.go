package resources

import (
	"backend/models"
	"time"
)

type AssociationResource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	Code        string `json:"code"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`

	Owner BasicUserResource `json:"user"`
}

func NewAssociationResource(association models.Association) AssociationResource {
	return AssociationResource{
		ID:          association.ID,
		Name:        association.Name,
		Description: association.Description,
		IsActive:    association.IsActive,
		Code:        association.Code,
		CreatedAt:   association.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   association.UpdatedAt.Format(time.RFC3339),
		Owner:       NewBasicUserResource(association.Owner),
	}
}
