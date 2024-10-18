package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	gorm.Model
	ID            string    `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" validate:"required,min=2,max=50"`
	Email         string    `gorm:"uniqueIndex:idx_email_deleted_at" json:"email" validate:"email,required"`
	Password      string    `json:"-"`
	PlainPassword *string   `gorm:"-" json:"password,omitempty" validate:"required_without=Password,omitempty,min=8,max=72"`
	Role          Role      `gorm:"default:user" json:"role" validate:"omitempty,oneof=admin user"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Associations   []Association   `json:"associations" gorm:"foreignKey:UserID"`
	Memberships    []Membership    `json:"memberships" gorm:"foreignKey:UserID"`
	Messages       []Message       `json:"messages" gorm:"foreignKey:SenderID"`
	Participations []Participation `json:"participations" gorm:"foreignKey:ParticipationID"`
	Events         []Event         `json:"events" gorm:"many2many:user_events;"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (u User) IsAdmin() bool {
	return u.Role == AdminRole
}
