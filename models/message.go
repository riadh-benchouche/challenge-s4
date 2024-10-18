package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageType string

type Message struct {
	gorm.Model
	ID            string      `json:"id" gorm:"primaryKey"`
	Content       string      `json:"content" validate:"required,min=10,max=300"`
	AssociationID uint        `json:"association_id" validate:"required"`
	Association   Association `gorm:"foreignkey:AssociationID" json:"association"`
	UserID        uint        `json:"user_id" validate:"required"`
	User          User        `gorm:"foreignkey:UserID" json:"user"`
	EventID       *uint       `json:"event_id"`
	Event         *Event      `gorm:"foreignkey:EventID" json:"event,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
}

type MessageCreate struct {
	Content       string `json:"content" validate:"required,min=10,max=300"`
	AssociationID uint   `json:"association_id" validate:"required"`
	UserID        uint   `json:"-"`
	EventID       *uint  `json:"event_id"`
}

type MessageUpdate struct {
	Content string `json:"content" validate:"required,min=10,max=300"`
}

func (e MessageCreate) ToMessage() *Message {
	return &Message{
		Content:       e.Content,
		AssociationID: e.AssociationID,
		UserID:        e.UserID,
		EventID:       e.EventID,
	}
}
