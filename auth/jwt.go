package utils

import (
	"backend/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user models.User) string {
	// Créer un token avec l'algorithme HS256 et les claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expire dans 24h
	})

	// Signer le token avec votre clé secrète
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "votre_clé_secrète_par_défaut" // Pour les tests uniquement
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return ""
	}

	return tokenString
}

// Fonction utile pour vérifier un token (optionnel pour les tests)
func VerifyJWT(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "votre_clé_secrète_par_défaut"
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
}
