package models

import (
	"time"

	"gorm.io/gorm"
)

type Association struct {
	gorm.Model
	ID          string    `json:"id" gorm:"primaryKey" validate:"required"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" default:"false"`
	Code        string    `json:"code" gorm:"unique;not null" validate:"required,min=5,max=20"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ImageURL    string    `json:"image_url"`

	// Foreign keys
	OwnerID string `json:"owner_id" validate:"required"`

	// Relationships
	Owner    User      `gorm:"foreignKey:OwnerID" json:"owner"`
	Members  []User    `gorm:"many2many:memberships;joinForeignKey:AssociationID;joinReferences:UserID" json:"members"`
	Messages []Message `gorm:"foreignKey:AssociationID"`
}

func (a Association) ToAssociation() *Association {
	return &Association{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Code:        a.Code,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		OwnerID:     a.OwnerID,
	}
}
