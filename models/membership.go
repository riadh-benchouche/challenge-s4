package models

import (
	"backend/enums"
	"gorm.io/gorm"
	"time"
)

type Membership struct {
	gorm.Model
	ID       string       `json:"id" gorm:"primaryKey"`
	JoinedAt time.Time    `json:"joined_at"`
	Status   enums.Status `json:"status" gorm:"default:pending" validate:"omitempty,oneof=pending accepted rejected"`
	Note     string       `json:"note"`

	// Foreign keys
	UserID        string `json:"user_id" validate:"required"`
	AssociationID string `json:"association_id" validate:"required"`

	// Relationships
	User        User        `gorm:"foreignKey:UserID" json:"user"`
	Association Association `gorm:"foreignKey:AssociationID" json:"association"`
}
