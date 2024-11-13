package models

import (
	"backend/enums"
	"time"

	"gorm.io/gorm"
)

type Participation struct {
	Status    enums.Status `json:"status" gorm:"default:pending" validate:"omitempty,oneof=pending present absent" faker:"oneof:pending,present,absent"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	// Foreign keys
	UserID  *string `json:"user_id" validate:"required" gorm:"primaryKey" faker:"-"`
	EventID *string `json:"event_id" gorm:"primaryKey" faker:"-"`

	// Relationships
	User  *User  `gorm:"foreignkey:UserID" json:"user" faker:"-"`
	Event *Event `gorm:"foreignkey:EventID" json:"event,omitempty" faker:"-"`
}

func (p *Participation) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = time.Now()
	return nil
}
