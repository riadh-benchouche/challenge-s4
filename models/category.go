package models

import (
	"backend/utils"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          string    `json:"id" gorm:"primaryKey" validate:"required"`
	Name        string    `json:"name" gorm:"not null" validate:"required" faker:"word"`
	Description string    `json:"description" faker:"sentence"`
	Note        int       `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = utils.GenerateULID()
	c.CreatedAt = time.Now()
	c.Note = utils.GenerateRandomNote()
	return nil
}
