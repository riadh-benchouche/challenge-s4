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

func (c *HomeController) GetStatistics(ctx echo.Context) error {
	loggedUser := ctx.Get("user").(models.User)
	stats, err := c.homeService.GetStatistics(loggedUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, stats)
}

func (c *HomeController) GetTopAssociations(ctx echo.Context) error {
	associations, err := c.homeService.GetTopAssociations()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, associations)
}
