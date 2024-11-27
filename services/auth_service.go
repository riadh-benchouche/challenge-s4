package services

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"backend/requests"
	"backend/resources"
	"backend/utils"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	Token string `json:"token"`
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

	if !s.CheckPasswordHash(password, targetUser.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	jwtSecret, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return nil, errors.ErrInternal
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    targetUser.ID,
			"name":  targetUser.Name,
			"email": targetUser.Email,
			"role":  targetUser.Role,
			"exp":   time.Now().Add(4 * time.Hour).Unix(),
			"iat":   time.Now().Unix(),
		},
	)

	token, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &LoginResponse{Token: token}, nil
}

func (s *AuthService) Register(Request requests.RegisterRequest) (*RegisterResponse, error) {
	hashedPassword, err := s.HashPassword(Request.Password)
	if err != nil {
		return nil, errors.ErrInternal
	}

	verificationToken := utils.GenerateULID()
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	// Log du token de vérification
	fmt.Printf("Creating user with verification token: %s\n", verificationToken)

	newUser := models.User{
		ID:                utils.GenerateULID(),
		Name:              Request.Name,
		Email:             Request.Email,
		Password:          hashedPassword,
		Role:              "user",
		IsActive:          false,
		VerificationToken: verificationToken,
		TokenExpiresAt:    &tokenExpiresAt,
	}

	if err := database.CurrentDatabase.Create(&newUser).Error; err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return nil, errors.ErrInternal
	}

	// Vérification après création
	var createdUser models.User
	database.CurrentDatabase.First(&createdUser, "id = ?", newUser.ID)
	fmt.Printf("Created user with token: %v\n", createdUser.VerificationToken)

	jwtSecret, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return nil, errors.ErrInternal
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    newUser.ID,
		"email": newUser.Email,
		"exp":   time.Now().Add(4 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}).SignedString([]byte(jwtSecret))

	if err != nil {
		return nil, errors.ErrInternal
	}

	userResource := resources.NewUserResource(newUser)
	return &RegisterResponse{User: userResource, Token: token}, nil
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

func (s *AuthService) VerifyEmail(token string) error {
	fmt.Printf("Attempting to verify token: %s\n", token)

	var user models.User
	result := database.CurrentDatabase.Where(
		"verification_token = ? AND email_verified_at IS NULL AND token_expires_at > ?",
		token,
		time.Now(),
	).First(&user)

	if result.Error != nil {
		fmt.Printf("Error finding user with token: %v\n", result.Error)
		return errors.ErrInvalidToken
	}

	fmt.Printf("Found user: %s with token: %s\n", user.Email, user.VerificationToken)

	now := time.Now()
	updates := map[string]interface{}{
		"email_verified_at":  now,
		"is_active":          true,
		"verification_token": "",
		"token_expires_at":   nil,
	}

	if err := database.CurrentDatabase.Model(&user).Updates(updates).Error; err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *AuthService) RegenerateVerificationToken(email string) error {
	var user models.User
	result := database.CurrentDatabase.Where(
		"email = ? AND email_verified_at IS NULL",
		email,
	).First(&user)

	if result.Error != nil {
		return errors.ErrUserNotFound
	}

	verificationToken := utils.GenerateULID()
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	updates := map[string]interface{}{
		"verification_token": verificationToken,
		"token_expires_at":   tokenExpiresAt,
	}

	if err := database.CurrentDatabase.Model(&user).Updates(updates).Error; err != nil {
		return errors.ErrInternal
	}

	return nil
}
