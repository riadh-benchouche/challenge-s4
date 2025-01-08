package test_utils

import (
	"backend/database"
	"backend/enums"
	"backend/models"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = database.InitTestDB()
	if err != nil {
		panic(fmt.Sprintf("Échec initialisation BD test: %v", err))
	}
}

func SetupTestDB() error {
	if err := CleanTestDB(); err != nil {
		return fmt.Errorf("échec nettoyage BD: %v", err)
	}

	err := db.AutoMigrate(database.Models...)
	if err != nil {
		return fmt.Errorf("échec migration BD test: %v", err)
	}

	database.CurrentDatabase = db
	return nil
}

func CleanTestDB() error {
	tables := []string{"participations", "events", "memberships", "associations", "users"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("échec suppression table %s: %v", table, err)
		}
	}
	return nil
}

func GetValidUser(role enums.Role) *models.User {
	now := time.Now()
	return &models.User{
		ID:              ulid.Make().String(),
		Name:            "Test User",
		Email:           fmt.Sprintf("test.%s@example.com", ulid.Make().String()),
		Role:            role,
		IsConfirmed:     true,
		IsActive:        true,
		EmailVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func GetAdminUser() *models.User {
	return GetValidUser(enums.AdminRole)
}

func GetAuthenticatedUser() *models.User {
	return GetValidUser(enums.UserRole)
}

func GetValidAssociation() *models.Association {
	return &models.Association{
		ID:          ulid.Make().String(),
		Name:        "Test Association",
		Description: "Test Description",
		IsActive:    false,
		ImageURL:    "https://test.com/image.jpg",
	}
}

func CreateUserAndAssociation() (*models.User, *models.Association) {
	user := GetAuthenticatedUser()
	if err := db.Create(user).Error; err != nil {
		panic(fmt.Sprintf("Échec création utilisateur: %v", err))
	}

	association := GetValidAssociation()
	association.OwnerID = user.ID
	association.Owner = *user

	if err := db.Create(association).Error; err != nil {
		panic(fmt.Sprintf("Échec création association: %v", err))
	}

	membership := &models.Membership{
		UserID:        user.ID,
		AssociationID: association.ID,
		JoinedAt:      time.Now(),
	}

	if err := db.Create(membership).Error; err != nil {
		panic(fmt.Sprintf("Échec création membership: %v", err))
	}

	return user, association
}

func GetValidEvent(associationID string) models.Event {
	return models.Event{
		ID:            ulid.Make().String(),
		Name:          "Test Event",
		Description:   "Test Description",
		Location:      "Test Location",
		AssociationID: associationID,
		Date:          time.Now().Add(24 * time.Hour),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func GetValidParticipation(userID string, eventID string) models.Participation {
	return models.Participation{
		ID:          ulid.Make().String(),
		UserID:      userID,
		EventID:     eventID,
		IsAttending: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
