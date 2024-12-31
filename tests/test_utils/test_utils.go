// backend/tests/test_utils/test_utils.go
package test_utils

import (
	"backend/database"
	"backend/enums"
	"backend/models"
	"backend/utils"
	"fmt"
	"time"
)

const Password = "TestPassword123!"

func SetupTestDB() error {
	db, err := database.InitTestDB()
	if err != nil {
		return fmt.Errorf("failed to initialize test database: %v", err)
	}

	// Nettoyer les tables
	if err := db.Exec("TRUNCATE TABLE associations CASCADE").Error; err != nil {
		return fmt.Errorf("failed to truncate associations: %v", err)
	}
	if err := db.Exec("TRUNCATE TABLE users CASCADE").Error; err != nil {
		return fmt.Errorf("failed to truncate users: %v", err)
	}

	return nil
}

// GetValidAssociation retourne une association valide pour les tests
func GetValidAssociation() models.Association {
	owner := GetAuthenticatedUser()
	return models.Association{
		ID:          utils.GenerateULID(),
		Name:        "Test Association",
		Description: "Test Description",
		Code:        utils.GenerateAssociationCode(),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerID:     owner.ID,
		Owner:       owner,
	}
}

// GetAuthenticatedUser retourne un utilisateur valide pour les tests
func GetAuthenticatedUser() models.User {
	plainPassword := Password
	return models.User{
		ID:            utils.GenerateULID(),
		Name:          "Test User",
		Email:         "test@example.com",
		PlainPassword: &plainPassword,
		Role:          enums.UserRole,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// CreateUserAndAssociation crée un utilisateur et une association pour les tests
func CreateUserAndAssociation() (models.User, models.Association) {
	SetupTestDB()
	user := GetAuthenticatedUser()
	database.CurrentDatabase.Create(&user)

	association := GetValidAssociation()
	association.OwnerID = user.ID
	association.Owner = user
	database.CurrentDatabase.Create(&association)
	return user, association
}

// CreateUser crée un utilisateur pour les tests
func CreateUser(role enums.Role) (*models.User, error) {
	user := GetAuthenticatedUser()
	user.Role = role
	err := database.CurrentDatabase.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserToken récupère un token pour un utilisateur donné
func GetUserToken(email string, password string) (*string, error) {
	user, err := CreateUser(enums.UserRole)
	if err != nil {
		return nil, err
	}
	token := generateTokenForUser(*user)
	return &token, nil
}

// GetTokenForNewUser crée un nouvel utilisateur et retourne son token
func GetTokenForNewUser(role enums.Role) (*string, error) {
	user, err := CreateUser(role)
	if err != nil {
		return nil, err
	}
	token := generateTokenForUser(*user)
	return &token, nil
}

// GetValidUser retourne un utilisateur valide avec le rôle spécifié
func GetValidUser(role enums.Role) models.User {
	user := GetAuthenticatedUser()
	user.Role = role
	return user
}

// Fonction helper privée pour générer un token
func generateTokenForUser(user models.User) string {
	return "test_token_" + user.ID // Pour les tests, on retourne un token simple
}
