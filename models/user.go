package models

import (
	"backend/enums"
	"backend/utils"
	"time"

	"gorm.io/gorm"
)

type TokenPair struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID              string     `json:"id" gorm:"primaryKey" validate:"required"`
	Name            string     `json:"name" validate:"required,min=2,max=50" faker:"name"`
	Email           string     `gorm:"uniqueIndex:idx_email_deleted_at" json:"email" validate:"email,required" faker:"email"`
	Password        string     `json:"-" faker:"password"`
	PlainPassword   *string    `gorm:"-" json:"plain_password,omitempty" validate:"required_without=Password,omitempty,min=8,max=20"`
	Role            enums.Role `gorm:"default:user" json:"role" validate:"omitempty,oneof=admin user association_leader" faker:"oneof:admin,association_leader,user"`
	IsConfirmed     bool       `json:"is_confirmed" gorm:"default:false"`
	IsActive        bool       `json:"is_active" gorm:"default:false" faker:"-"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ImageURL        string     `json:"image_url" faker:"url"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"index"`
	PointsOpen      int        `json:"points_open" gorm:"default:0"`
	FirebaseToken   string     `json:"firebase_token" validate:"omitempty"`

	AssociationsOwned []Association   `json:"associations_owned" gorm:"foreignKey:OwnerID" faker:"-"`
	Memberships       []Membership    `json:"memberships" gorm:"foreignKey:UserID" faker:"-"`
	Associations      []Association   `gorm:"many2many:memberships;joinForeignKey:UserID;joinReferences:AssociationID" json:"associations" faker:"-"`
	Messages          []Message       `json:"messages" gorm:"foreignKey:SenderID" faker:"-"`
	Participation     []Participation `json:"participation" gorm:"foreignKey:UserID" faker:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateULID()
	u.CreatedAt = time.Now()
	return nil
}

func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}
