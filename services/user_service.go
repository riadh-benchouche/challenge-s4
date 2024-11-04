package services

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"backend/utils"
	"strconv"

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

	create := database.CurrentDatabase.Create(&user)
	if create.Error != nil {
		return nil, create.Error
	}

	return &user, nil
}

func (s *UserService) GetUsers(pagination utils.Pagination, search *string) (*utils.Pagination, error) {
	var users []models.User
	query := database.CurrentDatabase

	if search != nil && *search != "" {
		searchedId, _ := strconv.Atoi(*search)

		query = query.Where(
			query.Where("id = ?", searchedId).
				Or("LOWER(name) LIKE ?", "%"+*search+"%").
				Or("LOWER(email) LIKE ?", "%"+*search+"%"))
	}

	err := query.Scopes(utils.Paginate(users, &pagination, query)).
		Order("ID asc").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	pagination.Rows = users

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
	result := database.CurrentDatabase.First(&user, "id = ?", id) // Rechercher l'utilisateur par ULID
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

	err = database.CurrentDatabase.Model(&models.User{}).Where("id = ?", id).Updates(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// func (s *UserService) GetUserEvents(userID uint, pagination utils.Pagination) (*utils.Pagination, error) {
// 	var events []models.Event

// 	query := database.CurrentDatabase.
// 		Joins("JOIN attends ON attends.event_id = events.id").
// 		Where("attends.user_id = ?", userID).
// 		Where("date >= ?", time.Now().Format(models.DateFormat)).
// 		Where("time is null or (date > ? or time >= ?)", time.Now().Format(models.DateFormat), time.Now().Format(models.TimeFormat)).
// 		Preload("Participants").
// 		Preload("Address").
// 		Order("date").
// 		Order("time")

// 	query.Scopes(utils.Paginate(events, &pagination, query)).Find(&events)

// 	pagination.Rows = events

// 	return &pagination, nil
// }
