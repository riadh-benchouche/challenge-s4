package services

import "backend/models"

type HomeService struct{}

func NewHomeService() *HomeService {
	return &HomeService{}
}

func (s *HomeService) GetHelloMessage() string {
	return "Welcome to our API!"
}

func (s *HomeService) GetHelloUserMessage(user models.User) string {
	return "Hello, " + user.Name + "! You are an logged in :D."
}
