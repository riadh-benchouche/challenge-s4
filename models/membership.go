package models

import (
	"backend/enums"
	"time"

	"gorm.io/gorm"
)

type Membership struct {
	gorm.Model
	JoinedAt time.Time    `json:"joined_at"`
	Status   enums.Status `json:"status" gorm:"default:pending" validate:"omitempty,oneof=pending accepted rejected" faker:"oneof:pending,accepted,rejected"`
	Note     string       `json:"note"`

	// Foreign keys
	UserID        string `json:"user_id" validate:"required" gorm:"primaryKey" faker:"-"`
	AssociationID string `json:"association_id" validate:"required" gorm:"primaryKey" faker:"-"`

	// Relationships
	User        User        `gorm:"foreignKey:UserID" json:"user"`
	Association Association `gorm:"foreignKey:AssociationID" json:"association"`
}
