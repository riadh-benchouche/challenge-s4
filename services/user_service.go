package services

import (
	"backend/config"
	"backend/database"
	"backend/enums"
	"backend/errors"
	"backend/models"
	"backend/resources"
	"backend/utils"
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) AddUser(user models.User) (*models.User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(user)
	if err != nil {
		return nil, err
	}

	var existingUser models.User
	database.CurrentDatabase.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != "" {
		return nil, errors.ErrUserAlreadyExists
	}

	user.Password, err = NewAuthService().HashPassword(*user.PlainPassword)
	if err != nil {
		return nil, err
	}
	user.PlainPassword = nil

	// Gérer la vérification d'email avec Redis
	verificationToken := utils.GenerateULID()
	ctx := context.Background()

	// Stocker le token dans Redis
	err = config.RedisClient.Set(ctx,
		"email_verification:"+verificationToken,
		user.Email,
		24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	user.IsActive = false
	user.EmailVerifiedAt = nil

	create := database.CurrentDatabase.Create(&user)
	if create.Error != nil {
		// En cas d'erreur, nettoyer Redis
		config.RedisClient.Del(ctx, "email_verification:"+verificationToken)
		return nil, create.Error
	}

	return &user, nil
}

func (s *UserService) GetUsers(pagination utils.Pagination, filters ...UserFilter) (*utils.Pagination, error) {
	var users []models.User
	query := database.CurrentDatabase.Model(&models.User{})

	if len(filters) > 0 {
		for _, filter := range filters {
			query = query.Where(filter.Column+" ILIKE ?", "%"+fmt.Sprintf("%v", filter.Value)+"%")
		}
	}

	err := query.Scopes(utils.Paginate(users, &pagination, query)).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	var userResources []resources.UserResource
	for _, user := range users {
		userResources = append(userResources, resources.NewUserResource(user))
	}

	pagination.Rows = userResources
	return &pagination, nil
}

func (s *UserService) DeleteUser(id string) error {
	var user models.User
	if err := database.CurrentDatabase.Where("id = ?", id).First(&user).Error; err != nil {
		return errors.ErrNotFound
	}
	if err := database.CurrentDatabase.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) FindByID(id string) (*models.User, error) {
	var user models.User
	result := database.CurrentDatabase.
		Preload("AssociationsOwned").
		Preload("Memberships").
		Preload("Associations").
		First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id string, user models.User) (*models.User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(user)
	if err != nil {
		return nil, err
	}

	var existingUser models.User
	database.CurrentDatabase.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != "" && existingUser.ID != id {
		return nil, errors.ErrUserAlreadyExists
	}

	if user.PlainPassword != nil {
		user.Password, err = NewAuthService().HashPassword(*user.PlainPassword)
		if err != nil {
			return nil, err
		}
		user.PlainPassword = nil
	}

	// Si l'email est modifié, réinitialiser la vérification
	var currentUser models.User
	if err := database.CurrentDatabase.First(&currentUser, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if user.Email != currentUser.Email {
		// Générer un nouveau token dans Redis
		verificationToken := utils.GenerateULID()
		ctx := context.Background()

		err = config.RedisClient.Set(ctx,
			"email_verification:"+verificationToken,
			user.Email,
			24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		user.EmailVerifiedAt = nil
		user.IsActive = false
	}

	err = database.CurrentDatabase.Model(&models.User{}).Where("id = ?", id).Updates(&user).Error
	if err != nil {
		return nil, err
	}

	var updatedUser models.User
	if err := database.CurrentDatabase.First(&updatedUser, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (s *UserService) JoinAssociation(userID, associationID, code string) (bool, error) {
	var user models.User
	if err := database.CurrentDatabase.First(&user, "id = ?", userID).Error; err != nil {
		return false, errors.ErrNotFound
	}

	var association models.Association
	if err := database.CurrentDatabase.First(&association, "id = ?", associationID).Error; err != nil {
		return false, errors.ErrNotFound
	}

	if association.Code != code {
		return false, errors.ErrInvalidCode
	}

	var membership models.Membership
	if err := database.CurrentDatabase.Where("user_id = ? AND association_id = ?", userID, associationID).First(&membership).Error; err == nil {
		return false, errors.ErrAlreadyJoined
	}

	newMembership := models.Membership{
		UserID:        userID,
		AssociationID: associationID,
		JoinedAt:      time.Now(),
		Status:        enums.Accepted,
	}

	if err := database.CurrentDatabase.Create(&newMembership).Error; err != nil {
		return false, err
	}

	return true, nil
}

type UserFilter struct {
	database.Filter
	Column string `json:"column" validate:"required,oneof=name email"`
}
