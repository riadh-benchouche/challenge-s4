package utils

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

// SendNotification envoie une notification push à un utilisateur
func SendNotification(firebaseToken, title, body string) error {
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile("firebase-service-account.json"))
	if err != nil {
		return fmt.Errorf("Erreur d'initialisation de Firebase Admin: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("Erreur d'initialisation du client Firebase Messaging: %v", err)
	}

	message := &messaging.Message{
		Token: firebaseToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'envoi de la notification: %v", err)
	}

	log.Printf("Notification envoyée avec succès : %s", response)
	return nil
}
