package controllers

import (
	// Remplacez par votre import path
	"backend/services" // Remplacez par votre import path
)

type MembershipController struct {
	service *services.MembershipService
}

func NewMembershipController(service *services.MembershipService) *MembershipController {
	return &MembershipController{service: service}
}
