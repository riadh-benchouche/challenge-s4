package controllers

import (
	"backend/database"
	coreErrors "backend/errors"
	"backend/models"
	"backend/requests"
	"backend/services"
	"backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (c *AuthController) Login(ctx echo.Context) error {
	var jsonBody requests.LoginRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(jsonBody)
	if err != nil {
		validationErrors := utils.GetValidationErrors(err.(validator.ValidationErrors), jsonBody)
		return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
	}

	result, err := c.authService.Login(jsonBody.Email, jsonBody.Password)
	if err != nil {
		if errors.Is(err, coreErrors.ErrInvalidCredentials) {
			return ctx.String(http.StatusUnauthorized, "Invalid credentials")
		}
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) Register(ctx echo.Context) error {
	var jsonBody requests.RegisterRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&jsonBody)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(jsonBody)
	if err != nil {
		validationErrors := utils.GetValidationErrors(err.(validator.ValidationErrors), jsonBody)
		return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
	}

	var existingUser models.User
	database.CurrentDatabase.Where("email = ?", jsonBody.Email).First(&existingUser)
	if existingUser.ID != "" {
		return ctx.String(http.StatusConflict, "Email already used")
	}

	result, err := c.authService.Register(jsonBody)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	confirmationLink := fmt.Sprintf("http://localhost:3000/confirm?token=%s", result.User.ID)
	subject := "Confirmation de votre inscription"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Confirmation d'inscription</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 20px;
            display: flex;
            justify-content: center;
        }
        .container {
            max-width: 600px;
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            text-align: center;
        }
        h2 {
            color: #333;
        }
        p {
            color: #555;
            line-height: 1.6;
        }
        .button {
            display: inline-block;
            background-color: #007bff;
            padding: 15px 25px;
            margin-top: 20px;
            color: #fff;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
        }
        .button:hover {
            background-color: #0056b3;
        }
        .footer {
            margin-top: 30px;
            font-size: 12px;
            color: #aaa;
        }
        .logo {
            max-width: 100px;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <img src="https://via.placeholder.com/100x100?text=Logo" alt="Logo" class="logo"/>
        <h2>Bienvenue, %s !</h2>
        <p>Merci de vous être inscrit. Pour activer votre compte, veuillez confirmer votre adresse email en cliquant sur le bouton ci-dessous :</p>
        <a href="%s" class="button">Confirmer mon inscription</a>
        <p>Ou copiez-collez ce lien dans votre navigateur :</p>
        <p><a href="%s">%s</a></p>
        <div class="footer">
            <p>Si vous n'avez pas initié cette inscription, vous pouvez ignorer cet email.</p>
            <p>&copy; 2024 VotreEntreprise. Tous droits réservés.</p>
        </div>
    </div>
</body>
</html>
`, result.User.Name, confirmationLink, confirmationLink, confirmationLink)

	// Envoyer un email de confirmation à l'utilisateur
	if err := utils.SendEmail(result.User.Email, subject, body); err != nil {
		ctx.Logger().Error("Erreur lors de l'envoi de l'email :", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to send confirmation email"})
	}

	return ctx.JSON(http.StatusCreated, result)
}
