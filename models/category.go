package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          string    `json:"id" gorm:"primaryKey" validate:"required"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Description string    `json:"description"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
