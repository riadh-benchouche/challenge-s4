package services

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"backend/requests"
	"backend/resources"
	"backend/utils"
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
	Token string      `json:"token"`
	User  models.User `json:"user"`
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

	return &LoginResponse{Token: token, User: targetUser}, nil
}

func (s *AuthService) Register(Request requests.RegisterRequest) (*RegisterResponse, error) {
	hashedPassword, err := s.HashPassword(Request.Password)
	if err != nil {
		return nil, errors.ErrInternal
	}

	newUser := models.User{
		ID:       utils.GenerateULID(),
		Name:     Request.Name,
		Email:    Request.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	database.CurrentDatabase.Create(&newUser)
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
