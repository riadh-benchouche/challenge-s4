package resources

import (
	"backend/models"
	"time"
)

type BasicUserResource struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserResource struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Email        string                `json:"email"`
	IsActive     bool                  `json:"is_active"`
	Role         string                `json:"role"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
	Associations []AssociationResource `json:"associations"`
	Memberships  []MembershipResource  `json:"memberships"`
}

func NewUserResource(user models.User) UserResource {
	associations := make([]AssociationResource, len(user.Associations))
	for i, association := range user.Associations {
		associations[i] = NewAssociationResource(association)
	}

	memberships := make([]MembershipResource, len(user.Memberships))
	for i, membership := range user.Memberships {
		memberships[i] = NewMembershipResource(membership)
	}

	return UserResource{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		IsActive:     user.IsActive,
		Role:         string(user.Role),
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		Associations: associations,
		Memberships:  memberships,
	}
}

func NewBasicUserResource(user models.User) BasicUserResource {
	return BasicUserResource{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
}
