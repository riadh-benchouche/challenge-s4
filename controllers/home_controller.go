package controllers

import (
	"backend/models"
	"backend/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HomeController struct {
	homeService *services.HomeService
}

func NewHomeController() *HomeController {
	return &HomeController{
		homeService: services.NewHomeService(),
	}
}

func (c *HomeController) Hello(ctx echo.Context) error {
	message := c.homeService.GetHelloMessage()
	return ctx.String(http.StatusOK, message)
}

func (c *HomeController) HelloAdmin(ctx echo.Context) error {
	loggedUser := ctx.Get("user").(models.User)
	message := c.homeService.GetHelloUserMessage(loggedUser)
	return ctx.String(http.StatusOK, message)
}
