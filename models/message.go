package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageType string

type Message struct {
	gorm.Model
	ID        string    `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" validate:"required,min=10,max=300"`
	CreatedAt time.Time `json:"created_at"`

	// Foreign keys
	AssociationID string `json:"association_id" validate:"required"`
	SenderID      string `json:"sender_id" validate:"required"`

	// Relationships
	Association Association `gorm:"foreignkey:AssociationID" json:"association"`
	Sender      User        `gorm:"foreignkey:SenderID" json:"user"`
}

type MessageCreate struct {
	Content       string `json:"content" validate:"required,min=1,max=300"`
	AssociationID string `json:"association_id" validate:"required"`
	SenderID      string `json:"sender_id" validate:"required"`
}

type MessageUpdate struct {
	Content string `json:"content" validate:"required,min=1,max=300"`
}

func (e MessageCreate) ToMessage() *Message {
	return &Message{
		Content:       e.Content,
		AssociationID: e.AssociationID,
		SenderID:      e.SenderID,
	}
}
