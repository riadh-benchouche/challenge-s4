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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Foreign keys
	OwnerID string `json:"owner_id" validate:"required"`

	// Relationships
	Owner       User         `gorm:"foreignkey:OwnerID" json:"user"`
	Memberships []Membership `gorm:"foreignKey:AssociationID"`
	Messages    []Message    `gorm:"foreignKey:AssociationID"`
}
