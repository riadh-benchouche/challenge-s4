package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)


type EmailConfig struct {
	Sender     string
	Identifier string
	Password   string
	Host       string
	Port       int
}

// getEmailConfig récupère la configuration email depuis les variables d'environnement
func getEmailConfig() (*EmailConfig, error) {
	config := &EmailConfig{
		Sender:     os.Getenv("EMAIL_SENDER"),
		Identifier: os.Getenv("EMAIL_IDENTIFIER"),
		Password:   os.Getenv("EMAIL_PASSWORD"),
		Host:       os.Getenv("SMTP_HOST"),
		Port:       587, // Port par défaut
	}

	// Vérification des variables d'environnement requises
	if config.Sender == "" || config.Password == "" || config.Host == "" || config.Identifier == "" {
		return nil, fmt.Errorf("SMTP configuration is missing: all environment variables must be set")
	}

	return config, nil
}

// SendEmail envoie un email avec les informations fournies
func SendEmail(to, subject, body string) error {
	config, err := getEmailConfig()
	if err != nil {
		fmt.Println("Erreur de configuration SMTP:", err)
		return err
	}

	fmt.Printf("📧 Envoi d'un email à %s\n", to)

	// Créer un nouveau message
	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configurer le serveur SMTP

	d := gomail.NewDialer(config.Host, config.Port, config.Identifier, config.Password)

	// Envoyer l'email
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf(" Erreur lors de l'envoi de l'email à %s: %v\n", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf(" Email envoyé avec succès à %s\n", to)
	return nil
}

// SendEmailWithRetry tente d'envoyer un email plusieurs fois en cas d'échec
func SendEmailWithRetry(to, subject, body string, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := SendEmail(to, subject, body); err != nil {
			lastErr = err
			fmt.Printf("Tentative %d/%d échouée: %v\n", i+1, maxRetries, err)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}

// ValidateEmail vérifie si l'adresse email est valide (à implémenter si nécessaire)
func ValidateEmail(email string) bool {
	// TODO: Implémenter la validation d'email si nécessaire
	return email != ""

}
