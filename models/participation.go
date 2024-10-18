package models

import (
	"time"

	"gorm.io/gorm"
)

type Participation struct {
	gorm.Model
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id" validate:"required"`
	User      *User     `gorm:"foreignkey:UserID" json:"user"`
	EventID   *uint     `json:"event_id"`
	Event     *Event    `gorm:"foreignkey:EventID" json:"event,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
