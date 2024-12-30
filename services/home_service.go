package services

import (
	"backend/database"
	"backend/enums"
	"backend/models"
)

type HomeService struct{}

func NewHomeService() *HomeService {
	return &HomeService{}
}

func (s *HomeService) GetHelloMessage() string {
	return "Welcome to our API V2!"
}

func (s *HomeService) GetHelloUserMessage(user models.User) string {
	return "Hello, " + user.Name + "! You are an logged in :D."
}

func (s *HomeService) GetStatistics(user models.User) (map[string]interface{}, error) {
	// 1. Récupérer les memberships de l'utilisateur avec les associations préchargées
	var memberships []models.Membership
	err := database.CurrentDatabase.
		Where("user_id = ? AND status = ?", user.ID, enums.Accepted).
		Preload("Association.Events").
		Preload("Association.Members").
		Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	// Initialiser les compteurs
	totalAssociations := len(memberships)
	var totalEvents int
	var totalUsers int

	// Calculer le nombre total d'événements et d'utilisateurs
	for _, membership := range memberships {
		// Compter les événements
		totalEvents += len(membership.Association.Events)

		// Compter les membres
		totalUsers += len(membership.Association.Members)
	}

	// Retourner les statistiques
	return map[string]interface{}{
		"total_associations": totalAssociations,
		"total_events":       totalEvents,
		"total_users":        totalUsers,
	}, nil
}

func (s *HomeService) GetTopAssociations() ([]models.Association, error) {
	var associations []models.Association

	// Récupérer les associations avec un comptage de leurs membres
	err := database.CurrentDatabase.
		Model(&models.Association{}).
		Select("associations.*, COUNT(DISTINCT memberships.user_id) as member_count").
		Joins("LEFT JOIN memberships ON memberships.association_id = associations.id").
		Where("memberships.status = ?", enums.Accepted).
		Group("associations.id").
		Order("member_count DESC").
		Limit(3).
		Preload("Owner").
		Find(&associations).Error

	if err != nil {
		return nil, err
	}

	return associations, nil
}
