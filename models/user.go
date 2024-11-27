package models

import (
	"backend/enums"
	"backend/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                string     `json:"id" gorm:"primaryKey" validate:"required"`
	Name              string     `json:"name" validate:"required,min=2,max=50" faker:"name"`
	Email             string     `gorm:"uniqueIndex:idx_email_deleted_at" json:"email" validate:"email,required" faker:"email"`
	Password          string     `json:"-" faker="password"`
	PlainPassword     *string    `gorm:"-" json:"password,omitempty" validate:"required_without=Password,omitempty,min=8,max=72"`
	Role              enums.Role `gorm:"default:user" json:"role" validate:"omitempty,oneof=admin user root" faker:"oneof:admin,association_leader,user"`
	IsActive          bool       `json:"is_active" gorm:"default:false" faker:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ImageURL          string     `json:"image_url" faker:"url"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at" gorm:"index"`                   // Index ajouté pour les recherches fréquentes
	VerificationToken string     `json:"verification_token,omitempty" gorm:"unique;index"` // Index ajouté pour les recherches par token
	TokenExpiresAt    *time.Time `json:"token_expires_at,omitempty" gorm:"index"`          // Index ajouté pour les vérifications d'expiration

	// Relationships
	AssociationsOwned []Association   `json:"associations_owned" gorm:"foreignKey:OwnerID" faker:"-"`
	Memberships       []Membership    `json:"memberships" gorm:"foreignKey:UserID" faker:"-"`
	Associations      []Association   `gorm:"many2many:memberships;joinForeignKey:UserID;joinReferences:AssociationID" json:"associations" faker:"-"`
	Messages          []Message       `json:"messages" gorm:"foreignKey:SenderID" faker:"-"`
	Participation     []Participation `json:"participation" gorm:"foreignKey:UserID" faker:"-"`
}

// BeforeCreate est appelé avant la création d'un utilisateur
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateULID()
	u.CreatedAt = time.Now()
	return nil
}

// IsEmailVerified vérifie si l'email est confirmé
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// IsTokenValid vérifie si le token de vérification est toujours valide
func (u *User) IsTokenValid() bool {
	return u.TokenExpiresAt != nil && u.TokenExpiresAt.After(time.Now())
}

// ClearVerificationData nettoie les données de vérification après confirmation
func (u *User) ClearVerificationData() {
	u.VerificationToken = ""
	u.TokenExpiresAt = nil
}
