package models

import (
	"backend/utils"
	"time"

	"gorm.io/gorm"
)

type Participation struct {
	ID          string    `json:"id" gorm:"primaryKey" validate:"required"`
	IsAttending bool      `json:"is_attending" gorm:"default:false" validate:"-" faker:"bool"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Foreign keys
	UserID  string `json:"user_id" validate:"required" gorm:"primaryKey" faker:"-"`
	EventID string `json:"event_id" gorm:"primaryKey" faker:"-"`

	// Relationships
	User  *User  `gorm:"foreignkey:UserID" json:"user" faker:"-"`
	Event *Event `gorm:"foreignkey:EventID" json:"event,omitempty" faker:"-"`
}

func (p *Participation) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = utils.GenerateULID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Participation) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return nil
}
