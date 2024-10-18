package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Foreign keys
	CategoryID    string `json:"category_id" `
	AssociationID string `json:"association_id" validate:"required"`

	// Relationships
	Category      Category        `gorm:"foreignkey:CategoryID" json:"category"`
	Association   Association     `gorm:"foreignkey:AssociationID" json:"association"`
	Participation []Participation `gorm:"foreignkey:EventID" json:"participation,omitempty"`
}
