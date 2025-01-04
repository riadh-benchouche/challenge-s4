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
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *AuthService) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type RegisterResponse struct {
	User  resources.UserResource `json:"user"`
	Token string                 `json:"token"`
}

func (s *AuthService) Login(email, password string) (*LoginResponse, error) {
	var targetUser models.User
	if err := database.CurrentDatabase.Where("email = ?", email).First(&targetUser).Error; err != nil {
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

	token, err := s.generateJWT(targetUser)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &LoginResponse{Token: token, User: targetUser}, nil
}

func (s *AuthService) Register(request requests.RegisterRequest) (*RegisterResponse, error) {
	ctx := context.Background()

	// Test de la connexion Redis
	if err := s.checkRedisConnection(ctx); err != nil {
		return nil, err
	}

	// Vérifier si l'email existe déjà
	if s.emailExists(request.Email) {
		return nil, errors.ErrEmailAlreadyExists
	}

	hashedPassword, err := s.HashPassword(request.Password)
	if err != nil {
		return nil, errors.ErrInternal
	}

	// Génération et stockage du token de vérification
	verificationToken := utils.GenerateULID()
	if err := s.storeVerificationToken(ctx, verificationToken, request.Email); err != nil {
		return nil, err
	}

	// Création de l'utilisateur
	newUser := models.User{
		ID:       utils.GenerateULID(),
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: false,
	}

	if err := database.CurrentDatabase.Create(&newUser).Error; err != nil {
		config.RedisClient.Del(ctx, fmt.Sprintf("email_verification:%s", verificationToken))
		return nil, errors.ErrInternal
	}

	// Génération du JWT
	token, err := s.generateJWT(newUser)
	if err != nil {
		return nil, err
	}

	userResource := resources.NewUserResource(newUser)
	userResource.VerificationToken = verificationToken
	return &RegisterResponse{User: userResource}, nil
}

func (s *AuthService) ConfirmEmail(token string) error {
	if token == "" {
		return errors.ErrInvalidToken
	}

	ctx := context.Background()
	email, err := s.getEmailFromToken(ctx, strings.TrimSpace(token))
	if err != nil {
		return err
	}

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
	config.RedisClient.Del(ctx, fmt.Sprintf("email_verification:%s", token))

	return nil
}

func (s *AuthService) ResendConfirmation(email string) error {
	if !s.unverifiedEmailExists(email) {
		return errors.ErrUserNotFound
	}

	verificationToken := utils.GenerateULID()
	ctx := context.Background()

	return s.storeVerificationToken(ctx, verificationToken, email)
}

// Méthodes utilitaires privées
func (s *AuthService) generateJWT(user models.User) (string, error) {
	jwtSecret, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return "", errors.ErrInternal
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
		"exp":   time.Now().Add(4 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *AuthService) checkRedisConnection(ctx context.Context) error {
	if _, err := config.RedisClient.Ping(ctx).Result(); err != nil {
		fmt.Printf("Redis connection error: %v\n", err)
		return errors.ErrInternal
	}
	return nil
}

func (s *AuthService) emailExists(email string) bool {
	var existingUser models.User
	return database.CurrentDatabase.Where("email = ?", email).First(&existingUser).Error == nil
}

func (s *AuthService) unverifiedEmailExists(email string) bool {
	var user models.User
	return database.CurrentDatabase.Where(
		"email = ? AND email_verified_at IS NULL",
		email,
	).First(&user).Error == nil
}

func (s *AuthService) storeVerificationToken(ctx context.Context, token, email string) error {
	key := fmt.Sprintf("email_verification:%s", token)
	if err := config.RedisClient.Set(ctx, key, email, 24*time.Hour).Err(); err != nil {
		fmt.Printf("Redis error setting token: %v\n", err)
		return errors.ErrInternal
	}
	return nil
}

func (s *AuthService) getEmailFromToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("email_verification:%s", token)
	email, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.ErrInvalidToken
	} else if err != nil {
		return "", errors.ErrInternal
	}
	return email, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	jwtSecret, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return nil, errors.ErrInternal
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
}
