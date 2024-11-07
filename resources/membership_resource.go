package resources

import (
	"backend/models"
	"time"
)

type MembershipResource struct {
	ID            string              `json:"id"`
	JoinedAt      string              `json:"joined_at"`
	Status        string              `json:"status"`
	Note          string              `json:"note"`
	UserID        string              `json:"user_id"`
	User          BasicUserResource   `json:"user"`
	AssociationID string              `json:"association_id"`
	Association   AssociationResource `json:"association"`
}

func NewMembershipResource(membership models.Membership) MembershipResource {
	return MembershipResource{
		ID:            membership.ID,
		JoinedAt:      membership.JoinedAt.Format(time.RFC3339),
		Status:        string(membership.Status),
		Note:          membership.Note,
		UserID:        membership.UserID,
		AssociationID: membership.AssociationID,
		User:          NewBasicUserResource(membership.User),
		Association:   NewAssociationResource(membership.Association),
	}
}
