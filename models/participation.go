package models

import (
	"backend/enum"
	"gorm.io/gorm"
	"time"
)

type Participation struct {
	gorm.Model
	ID        string      `json:"id" gorm:"primaryKey"`
	Status    enum.Status `json:"status" gorm:"default:pending" validate:"omitempty,oneof=pending present absent"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	// Foreign keys
	UserID  *string `json:"user_id" validate:"required"`
	EventID *string `json:"event_id"`

	// Relationships
	User  *User  `gorm:"foreignkey:UserID" json:"user"`
	Event *Event `gorm:"foreignkey:EventID" json:"event,omitempty"`
}
