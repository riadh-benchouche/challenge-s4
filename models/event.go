package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID              string           `json:"id" gorm:"primaryKey"`
	CategoryID      string          `json:"category_id" `
	Category        Category        `gorm:"foreignkey:CategoryID" json:"category"`
	Name            string          `json:"name" gorm:"not null"`
	Description     string          `json:"description"`
	AssociationID   string           `json:"association_id" validate:"required"`
	Association     Association     `gorm:"foreignkey:AssociationID" json:"association"`
	Date            time.Time       `json:"date"`
	Location        string          `json:"location"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	ParticipationID string          `json:"participationID"`
	Participation   []Participation `gorm:"foreignkey:ParicipationID" json:"participation"`
}
