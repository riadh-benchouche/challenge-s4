package models

import (
	"backend/utils"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null" faker:"word"`
	Description string    `json:"description" faker:"sentence"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location" faker:"word"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Foreign keys
	CategoryID    string `json:"category_id" faker:"-"`
	AssociationID string `json:"association_id" validate:"required" faker:"-"`

	// Relationships
	Category       Category        `gorm:"foreignKey:CategoryID" json:"category" validate:"-" faker:"-"`
	Association    Association     `gorm:"foreignKey:AssociationID" json:"association" validate:"-" faker:"-"`
	Participations []Participation `gorm:"foreignKey:EventID" json:"participations,omitempty" faker:"-"`
	User           []User          `gorm:"many2many:participations;joinForeignKey:EventID;joinReferences:UserID" json:"users" faker:"-"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = utils.GenerateULID()
	e.CreatedAt = time.Now()
	e.Date = time.Now()
	return nil
}
