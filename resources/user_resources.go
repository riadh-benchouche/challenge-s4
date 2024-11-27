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
	ID                string                `json:"id"`
	Name              string                `json:"name"`
	Email             string                `json:"email"`
	IsActive          bool                  `json:"is_active"`
	Role              string                `json:"role"`
	CreatedAt         string                `json:"created_at"`
	UpdatedAt         string                `json:"updated_at"`
	Associations      []AssociationResource `json:"associations"`
	Memberships       []MembershipResource  `json:"memberships"`
	VerificationToken string                `json:"verification_token,omitempty"`
	EmailVerifiedAt   *time.Time            `json:"email_verified_at"`
	ImageURL          string                `json:"image_url"`
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
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		IsActive:          user.IsActive,
		Role:              string(user.Role),
		CreatedAt:         user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         user.UpdatedAt.Format(time.RFC3339),
		Associations:      associations,
		Memberships:       memberships,
		ImageURL:          user.ImageURL,
		EmailVerifiedAt:   user.EmailVerifiedAt,
		VerificationToken: user.VerificationToken,
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

type ParticipationResource struct {
	IsAttending bool              `json:"is_attending"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	User        BasicUserResource `json:"user"`
	EventID     string            `json:"event_id"`
}

func NewParticipationResource(participation models.Participation) ParticipationResource {
	return ParticipationResource{
		IsAttending: participation.IsAttending,
		CreatedAt:   participation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   participation.UpdatedAt.Format(time.RFC3339),
		User:        NewBasicUserResource(*participation.User), // Utilisation de BasicUserResource
		EventID:     participation.EventID,
	}
}

func NewParticipationResourceList(participations []models.Participation) []ParticipationResource {
	resources := make([]ParticipationResource, len(participations))
	for i, participation := range participations {
		resources[i] = NewParticipationResource(participation)
	}
	return resources
}
