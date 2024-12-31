package models

import (
	"backend/utils"
	"time"

	"gorm.io/gorm"
)

type Association struct {
	ID          string    `json:"id" gorm:"primaryKey" validate:"required"`
	Name        string    `json:"name" gorm:"not null" faker:"name"`
	Description string    `json:"description" faker:"sentence"`
	IsActive    bool      `json:"is_active" gorm:"default:false"`
	Code        string    `json:"code" gorm:"unique;not null" validate:"required,min=5,max=20"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ImageURL    string    `json:"image_url" faker:"url"`

	// Foreign keys
	OwnerID string `json:"owner_id" validate:"required" faker:"-"`

	// Relationships
	Owner    User      `gorm:"foreignKey:OwnerID" json:"owner" faker:"-"`
	Members  []User    `gorm:"many2many:memberships;joinForeignKey:AssociationID;joinReferences:UserID" json:"members" faker:"-"`
	Messages []Message `gorm:"foreignKey:AssociationID" faker:"-"`
	Events   []Event   `gorm:"foreignKey:AssociationID" faker:"-"`
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

func (a *Association) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = utils.GenerateULID()
	a.Code = utils.GenerateAssociationCode()
	a.CreatedAt = time.Now()
	return nil
}
