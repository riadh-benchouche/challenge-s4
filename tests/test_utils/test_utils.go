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

	// Liste des tables à nettoyer
	tables := []string{
		"messages",
		"participations",
		"memberships",
		"events",
		"associations",
		"categories",
		"users",
	}

	// Désactiver les contraintes de clé étrangère temporairement
	if err := db.Exec("SET CONSTRAINTS ALL DEFERRED").Error; err != nil {
		return fmt.Errorf("failed to defer constraints: %v", err)
	}

	// Nettoyer toutes les tables
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %v", table, err)
		}
	}

	// Réactiver les contraintes
	if err := db.Exec("SET CONSTRAINTS ALL IMMEDIATE").Error; err != nil {
		return fmt.Errorf("failed to restore constraints: %v", err)
	}

	return nil
}

// GetValidUser retourne un utilisateur valide avec le rôle spécifié
func GetValidUser(role enums.Role) models.User {
	plainPassword := Password
	return models.User{
		ID:            utils.GenerateULID(),
		Name:          "Test User",
		Email:         fmt.Sprintf("test_%s@example.com", utils.GenerateULID()),
		PlainPassword: &plainPassword,
		Role:          role,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// GetAuthenticatedUser retourne un utilisateur standard pour les tests
func GetAuthenticatedUser() models.User {
	return GetValidUser(enums.UserRole)
}

// GetAdminUser retourne un utilisateur admin pour les tests
func GetAdminUser() models.User {
	return GetValidUser(enums.AdminRole)
}

// GetValidAssociation retourne une association valide pour les tests
func GetValidAssociation() models.Association {
	owner := GetAuthenticatedUser()
	return models.Association{
		ID:          utils.GenerateULID(),
		Name:        "Test Association",
		Description: "Test Description",
		IsActive:    false,
		Code:        utils.GenerateAssociationCode(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ImageURL:    "https://test.com/image.jpg",
		OwnerID:     owner.ID,
		Owner:       owner,
	}
}

// GetValidCategory retourne une catégorie valide pour les tests
func GetValidCategory() models.Category {
	return models.Category{
		ID:          utils.GenerateULID(),
		Name:        "Test Category",
		Description: "Test Category Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// GetValidEvent retourne un événement valide pour les tests
func GetValidEvent(associationID string) models.Event {
	return models.Event{
		ID:            utils.GenerateULID(),
		Name:          "Test Event",
		Description:   "Test Description",
		Date:          time.Now().Add(24 * time.Hour),
		Location:      "Test Location",
		AssociationID: associationID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// GetValidMembership retourne une adhésion valide pour les tests
func GetValidMembership(userID, associationID string) models.Membership {
	return models.Membership{
		JoinedAt:      time.Now(),
		Status:        enums.Pending,
		Note:          "Test Note",
		UserID:        userID,
		AssociationID: associationID,
	}
}

// GetValidParticipation retourne une participation valide pour les tests
func GetValidParticipation(userID, eventID string) models.Participation {
	return models.Participation{
		Status:    enums.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    &userID,
		EventID:   &eventID,
	}
}

// GetValidMessage retourne un message valide pour les tests
func GetValidMessage(senderID, associationID string) models.Message {
	return models.Message{
		ID:            utils.GenerateULID(),
		Content:       "Test message content with minimum length required",
		CreatedAt:     time.Now(),
		AssociationID: associationID,
		SenderID:      senderID,
	}
}

// CreateUser crée un utilisateur pour les tests
func CreateUser(role enums.Role) (*models.User, error) {
	user := GetValidUser(role)
	err := database.CurrentDatabase.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
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

// CreateTestData crée un jeu complet de données de test
func CreateTestData() (models.User, models.Association, models.Event, error) {
	if err := SetupTestDB(); err != nil {
		return models.User{}, models.Association{}, models.Event{}, err
	}

	user := GetAuthenticatedUser()
	if err := database.CurrentDatabase.Create(&user).Error; err != nil {
		return models.User{}, models.Association{}, models.Event{}, err
	}

	association := GetValidAssociation()
	association.OwnerID = user.ID
	if err := database.CurrentDatabase.Create(&association).Error; err != nil {
		return models.User{}, models.Association{}, models.Event{}, err
	}

	event := GetValidEvent(association.ID)
	if err := database.CurrentDatabase.Create(&event).Error; err != nil {
		return models.User{}, models.Association{}, models.Event{}, err
	}

	return user, association, event, nil
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

// generateTokenForUser génère un token de test pour un utilisateur
func generateTokenForUser(user models.User) string {
	return "test_token_" + user.ID
}
