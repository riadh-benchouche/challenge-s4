package models

import (
	"time"

	"gorm.io/gorm"
)

type Association struct {
	gorm.Model
	ID          string       `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"not null"`
	Description string       `json:"description"`
	UserID      uint         `json:"user_id" validate:"required"`
	User        User         `gorm:"foreignkey:UserID" json:"user"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Events      []Event      `gorm:"foreignKey:EventID"`
	Messages    []Message    `gorm:"foreignKey:MessageID"`
	Memberships []Membership `gorm:"foreignKey:MembershipID"`
}
