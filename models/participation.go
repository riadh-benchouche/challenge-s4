package models

import (
	"backend/enums"
	"time"

	"gorm.io/gorm"
)

type Participation struct {
	gorm.Model
	Status    enums.Status `json:"status" gorm:"default:pending" validate:"omitempty,oneof=pending present absent"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	// Foreign keys
	UserID  *string `json:"user_id" validate:"required" gorm:"primaryKey"`
	EventID *string `json:"event_id" gorm:"primaryKey"`

	// Relationships
	User  *User  `gorm:"foreignkey:UserID" json:"user"`
	Event *Event `gorm:"foreignkey:EventID" json:"event,omitempty"`
}
