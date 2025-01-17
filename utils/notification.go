package utils

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
)

// SendNotification envoie une notification push à un utilisateur
func SendNotification(firebaseToken, title, body string) error {
	fmt.Printf("Tentative d'envoi de notification - Token: %s\n", firebaseToken)

	// Créer la configuration Firebase
	config := &firebase.Config{
		ProjectID: "challengeflutter-5d5eb", // Utilisez l'ID de votre projet Firebase
	}

	// Initialiser l'app avec la configuration
	app, err := firebase.NewApp(context.Background(), config, option.WithCredentialsFile("firebase-service-account.json"))
	if err != nil {
		fmt.Printf("Erreur lors de l'initialisation de Firebase: %v\n", err)
		return fmt.Errorf("Erreur d'initialisation de Firebase Admin: %v", err)
	}
	fmt.Println("Firebase App initialisée avec succès")

	client, err := app.Messaging(context.Background())
	if err != nil {
		fmt.Printf("Erreur lors de l'initialisation du client Messaging: %v\n", err)
		return fmt.Errorf("Erreur d'initialisation du client Firebase Messaging: %v", err)
	}
	fmt.Println("Client Messaging initialisé avec succès")

	message := &messaging.Message{
		Token: firebaseToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		fmt.Printf("Erreur lors de l'envoi de la notification: %v\n", err)
		return fmt.Errorf("Erreur lors de l'envoi de la notification: %v", err)
	}

	fmt.Printf("Notification envoyée avec succès, response: %s\n", response)
	return nil
}
