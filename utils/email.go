package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail envoie un email avec les informations fournies
func SendEmail(to, subject, body string) error {
	sender := os.Getenv("EMAIL_SENDER")
	identifier := os.Getenv("EMAIL_IDENTIFIER")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587

	fmt.Println("Envoi d'un email à", to)
	// Vérification des variables d'environnement
	if sender == "" || password == "" || smtpHost == "" || identifier == "" {
		fmt.Println("Erreur : Les informations SMTP ne sont pas définies correctement.")
		return fmt.Errorf("SMTP configuration is missing")
	}

	// Créer un nouveau message
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configurer le serveur SMTP
	d := gomail.NewDialer(smtpHost, smtpPort, identifier, password)

	// Envoyer l'email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Erreur lors de l'envoi de l'email :", err)
		return err
	}

	fmt.Println("Email envoyé avec succès à", to)
	return nil

}
