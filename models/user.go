package models

import (
	"backend/enums"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            string     `json:"id" gorm:"primaryKey" validate:"required"`
	Name          string     `json:"name" validate:"required,min=2,max=50"`
	Email         string     `gorm:"uniqueIndex:idx_email_deleted_at" json:"email" validate:"email,required"`
	Password      string     `json:"-"`
	PlainPassword *string    `gorm:"-" json:"password,omitempty" validate:"required_without=Password,omitempty,min=8,max=72"`
	Role          enums.Role `gorm:"default:user" json:"role" validate:"omitempty,oneof=admin user root"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ImageURL      string     `json:"image_url"`

	// Relationships
	AssociationsOwned []Association   `json:"associations_owned" gorm:"foreignKey:OwnerID"`
	Memberships       []Membership    `json:"memberships" gorm:"foreignKey:UserID"`
	Associations      []Association   `gorm:"many2many:memberships;joinForeignKey:UserID;joinReferences:AssociationID" json:"associations"`
	Messages          []Message       `json:"messages" gorm:"foreignKey:SenderID"`
	Participation     []Participation `json:"participation" gorm:"foreignKey:UserID"`
}
