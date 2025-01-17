package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"backend/config"
	"backend/database"
	"backend/errors"
	"backend/models"
	"backend/requests"
	"backend/resources"
	"backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

type RegisterResponse struct {
	User  resources.UserResource `json:"user"`
	Token string                 `json:"token"`
}

func (s *AuthService) Login(email, password string) (*LoginResponse, error) {
	var targetUser models.User
	database.CurrentDatabase.Where("email = ?", email).First(&targetUser)

	if targetUser.ID == "" {
		return nil, errors.ErrInvalidCredentials
	}

	if !targetUser.IsActive || !targetUser.IsConfirmed {
		return nil, errors.ErrUserNotActive
	}

	if !s.CheckPasswordHash(password, targetUser.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	if !targetUser.IsActive || targetUser.EmailVerifiedAt == nil {
		return nil, errors.ErrEmailNotVerified
	}

	jwtSecret, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return nil, errors.ErrInternal
	}

	// Générer l'access token (courte durée)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    targetUser.ID,
			"name":  targetUser.Name,
			"email": targetUser.Email,
			"role":  targetUser.Role,
			"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(), // 15 minutes
			"iat":   time.Now().Unix(),
		},
	)

	// Générer le refresh token (longue durée)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  targetUser.ID,
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 jours
			"iat": time.Now().Unix(),
		},
	)

	// Signer les tokens
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.ErrInternal
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.ErrInternal
	}

	// Stocker le refresh token dans Redis
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", targetUser.ID)
	err = config.RedisClient.Set(ctx, key, refreshTokenString, 7*24*time.Hour).Err()
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &LoginResponse{
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
		User:         targetUser,
	}, nil
}

func (s *AuthService) Register(request requests.RegisterRequest) (*RegisterResponse, error) {
	ctx := context.Background()

	// Test de la connexion Redis
	pong, err := config.RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf(" Redis connection error: %v\n", err)
		return nil, errors.ErrInternal
	}
	fmt.Printf("🔄 Redis connection test: %v\n", pong)

	// Vérifier si l'email existe déjà
	var existingUser models.User
	if result := database.CurrentDatabase.Where("email = ?", request.Email).First(&existingUser); result.Error == nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	hashedPassword, err := s.HashPassword(request.Password)
	if err != nil {
		return nil, errors.ErrInternal
	}

	// Génération du token de vérification
	verificationToken := utils.GenerateULID()
	key := fmt.Sprintf("email_verification:%s", verificationToken)

	// Stockage dans Redis (expire après 24h)
	err = config.RedisClient.Set(ctx, key, request.Email, 24*time.Hour).Err()
	if err != nil {
		fmt.Printf(" Redis error setting token: %v\n", err)
		return nil, errors.ErrInternal
	}

	newUser := models.User{
		ID:            utils.GenerateULID(),
		Name:          request.Name,
		Email:         request.Email,
		Password:      hashedPassword,
		FirebaseToken: request.FirebaseToken,
		Role:          "user",
		IsActive:      false,
	}

	if err := database.CurrentDatabase.Create(&newUser).Error; err != nil {
		config.RedisClient.Del(ctx, key)
		return nil, errors.ErrInternal
	}

	userResource := resources.NewUserResource(newUser)
	userResource.VerificationToken = verificationToken
	return &RegisterResponse{User: userResource}, nil
}

func (s *AuthService) ConfirmEmail(token string) error {
	ctx := context.Background()

	if token == "" {
		fmt.Printf(" Token is empty\n")
		return errors.ErrInvalidToken
	}

	token = strings.TrimSpace(token)
	key := fmt.Sprintf("email_verification:%s", token)

	// Récupération de l'email associé au token
	email, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Printf(" Token not found in Redis: %s\n", token)
		return errors.ErrInvalidToken
	} else if err != nil {
		fmt.Printf(" Redis error: %v\n", err)
		return errors.ErrInternal
	}

	// Mise à jour de l'utilisateur
	now := time.Now()
	result := database.CurrentDatabase.Model(&models.User{}).
		Where("email = ? AND email_verified_at IS NULL", email).
		Updates(map[string]interface{}{
			"email_verified_at": now,
			"is_confirmed":      true,
			"is_active":         true,
		})

	if result.Error != nil {
		return errors.ErrInternal
	}

	if result.RowsAffected == 0 {
		return errors.ErrInvalidToken
	}

	// Suppression du token utilisé
	config.RedisClient.Del(ctx, key)

	return nil
}

func (s *AuthService) ResendConfirmation(email string) error {
	ctx := context.Background()

	var user models.User
	result := database.CurrentDatabase.Where(
		"email = ? AND email_verified_at IS NULL",
		email,
	).First(&user)

	if result.Error != nil {
		return errors.ErrUserNotFound
	}

	verificationToken := utils.GenerateULID()
	key := fmt.Sprintf("email_verification:%s", verificationToken)

	err := config.RedisClient.Set(ctx, key, email, 24*time.Hour).Err()
	if err != nil {
		return errors.ErrInternal
	}

	// Le token est généré et stocké, il devra être envoyé par email
	return nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	jwtSecret, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return nil, errors.ErrInternal
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}

func (s *AuthService) GenerateTokenPair(user models.User) (*models.TokenPair, error) {
	// Générer l'access token (comme avant, mais avec une durée plus courte)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
			"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(), // Durée plus courte
			"iat":   time.Now().Unix(),
		})

	// Générer le refresh token (plus long)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 jours
			"iat": time.Now().Unix(),
		})

	jwtSecret := os.Getenv("JWT_KEY")

	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	// Stocker le refresh token dans Redis
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", user.ID)
	err = config.RedisClient.Set(ctx, key, refreshTokenString, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenPair, error) {
	// Valider le refresh token
	token, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrInvalidToken
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return nil, errors.ErrInvalidToken
	}

	// Vérifier si le refresh token est dans Redis
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", userID)
	storedToken, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil || storedToken != refreshToken {
		return nil, errors.ErrInvalidToken
	}

	// Récupérer l'utilisateur
	var user models.User
	if err := database.CurrentDatabase.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.ErrUserNotFound
	}

	// Générer une nouvelle paire de tokens
	return s.GenerateTokenPair(user)
}
