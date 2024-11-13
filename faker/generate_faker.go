package faker

import (
	"backend/models"
	"backend/utils"
	"log"
	"time"

	"github.com/bxcodec/faker/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateFakeUser() models.User {
	var user models.User
	err := faker.FakeData(&user)
	if err != nil {
		log.Fatalf("Erreur de génération de données pour User : %v", err)
	}

	rawPassword := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 6)
	if err != nil {
		log.Fatalf("Erreur de hachage du mot de passe : %v", err)
	}

	user.ID = utils.GenerateULID()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user
}

func GenerateFakeEvent(categoryID, associationID string) models.Event {
	var event models.Event
	err := faker.FakeData(&event)
	if err != nil {
		log.Fatalf("Erreur de génération de données pour Event : %v", err)
	}
	event.ID = utils.GenerateULID()
	event.CategoryID = categoryID
	event.AssociationID = associationID
	event.CreatedAt = time.Now()
	return event
}

func GenerateFakeMessage(associationID, senderID string) models.Message {
	var message models.Message
	err := faker.FakeData(&message)
	if err != nil {
		log.Fatalf("Erreur de génération de données pour Message : %v", err)
	}
	message.ID = utils.GenerateULID()
	message.AssociationID = associationID
	message.SenderID = senderID
	message.CreatedAt = time.Now()
	return message
}

// func GenerateFakeParticipation() models.Participation {
// 	var participation models.Participation
// 	err := faker.FakeData(&participation)
// 	if err != nil {
// 		log.Fatalf("Erreur de génération de données pour Participation : %v", err)
// 	}
// 	participation.CreatedAt = time.Now()
// 	return participation
// }

func GenerateFakeAssociation(ownerID string) models.Association {
	var association models.Association
	err := faker.FakeData(&association)
	if err != nil {
		log.Fatalf("Erreur de génération de données pour Association : %v", err)
	}
	association.ID = utils.GenerateULID()
	association.OwnerID = ownerID
	association.Code = utils.GenerateAssociationCode()
	association.CreatedAt = time.Now()
	return association
}

func GenerateFakeCategory() models.Category {
	var category models.Category
	err := faker.FakeData(&category)
	if err != nil {
		log.Fatalf("Erreur de génération de données pour Category : %v", err)
	}
	category.ID = utils.GenerateULID()
	category.CreatedAt = time.Now()
	return category
}

// func GenerateFakeMembership() models.Membership {
// 	var membership models.Membership
// 	err := faker.FakeData(&membership)
// 	if err != nil {
// 		log.Fatalf("Erreur de génération de données pour Membership : %v", err)
// 	}
// 	membership.CreatedAt = time.Now()
// 	return membership
// }

func GenerateFakeData(db *gorm.DB) {

	var categories []models.Category
	var users []models.User
	for i := 0; i < 10; i++ {
		category := GenerateFakeCategory()
		db.Create(&category)
		categories = append(categories, category)
	}

	for i := 0; i < 10; i++ {
		user := GenerateFakeUser()
		db.Create(&user)
		users = append(users, user)
	}

	var associations []models.Association
	for i := 0; i < 10; i++ {
		ownerID := users[i%len(users)].ID
		association := GenerateFakeAssociation(ownerID)
		db.Create(&association)
		associations = append(associations, association)
	}

	for i := 0; i < 10; i++ {
		categoryID := categories[i%len(categories)].ID
		associationID := associations[i%len(associations)].ID
		event := GenerateFakeEvent(categoryID, associationID)
		db.Create(&event)
	}

	for i := 0; i < 10; i++ {
		senderID := users[i%len(users)].ID
		associationID := associations[i%len(associations)].ID
		message := GenerateFakeMessage(associationID, senderID)
		db.Create(&message)
	}

	// Étapes 5 et 6 : Générez des participations et memberships, si nécessaire, en suivant la même logique

	// for i := 0; i < 10; i++ {
	// 	participation := GenerateFakeParticipation()
	// 	db.Create(&participation)
	// }

	// for i := 0; i < 10; i++ {
	// 	membership := GenerateFakeMembership()
	// 	db.Create(&membership)
	// }
}
