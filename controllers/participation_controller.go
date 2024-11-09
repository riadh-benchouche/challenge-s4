package controllers

import (
	"backend/services"
)

type ParticipationController struct {
	service *services.ParticipationService
}

func NewParticipationController(service *services.ParticipationService) *ParticipationController {
	return &ParticipationController{service: service}
}
