package models

import (
	"time"

	"gorm.io/gorm"
)

type Membership struct {
	gorm.Model
	ID            string    `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id"`
	User          User      `gorm:"foreignkey:UserID" json:"user"`
	AssociationID uint      `json:"association_id"`
	JoinedAt      time.Time `json:"joined_at"`
	Note          string    `json:"note"`
}
