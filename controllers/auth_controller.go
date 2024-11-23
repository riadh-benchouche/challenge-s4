package controllers

import (
	"backend/database"
	"backend/errors"
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
		if err == errors.ErrInvalidCredentials {
			return ctx.String(http.StatusUnauthorized, "Invalid credentials")
		}
		if err == errors.ErrEmailNotVerified {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Email not verified",
			})
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

	result, err := c.authService.Register(jsonBody)
	if err != nil {
		if err == errors.ErrEmailAlreadyExists {
			return ctx.String(http.StatusConflict, "Email already used")
		}
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

func (c *AuthController) ConfirmEmail(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	if token == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Token is required",
		})
	}

	err := c.authService.ConfirmEmail(token)
	if err != nil {
		switch err {
		case errors.ErrInvalidToken:
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid or expired token",
			})
		default:
			ctx.Logger().Error(err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "An error occurred while confirming email",
			})
		}
	}

	// Récupérer l'utilisateur après la confirmation
	var user models.User
	if err := database.CurrentDatabase.Where("email_verified_at IS NOT NULL").
		Order("email_verified_at DESC").First(&user).Error; err != nil {
		ctx.Logger().Error("Erreur lors de la récupération de l'utilisateur :", err)
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "Email confirmé avec succès",
		})
	}

	subject := "Votre compte a été confirmé"
	body := fmt.Sprintf(`
   <!DOCTYPE html>
   <html lang="fr">
   <head>
       <meta charset="UTF-8">
       <title>Compte confirmé</title>
       <style>
           body {
               font-family: Arial, sans-serif;
               line-height: 1.6;
               color: #333;
           }
           .container {
               max-width: 600px;
               margin: 0 auto;
               padding: 20px;
               background-color: #f9f9f9;
               border-radius: 5px;
           }
           .header {
               text-align: center;
               color: #2c3e50;
           }
           .content {
               margin: 20px 0;
               padding: 20px;
               background-color: #ffffff;
               border-radius: 5px;
           }
       </style>
   </head>
   <body>
       <div class="container">
           <div class="header">
               <h2>Félicitations %s !</h2>
           </div>
           <div class="content">
               <p>Votre compte a été confirmé avec succès.</p>
               <p>Vous pouvez maintenant vous connecter et accéder à toutes les fonctionnalités de notre plateforme.</p>
               <p>Si vous avez des questions, n'hésitez pas à nous contacter.</p>
           </div>
       </div>
   </body>
   </html>
   `, user.Name)

	if err := utils.SendEmail(user.Email, subject, body); err != nil {
		ctx.Logger().Error("Erreur lors de l'envoi de l'email de confirmation :", err)
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Email confirmé avec succès",
	})
}

func (c *AuthController) ResendConfirmation(ctx echo.Context) error {
	email := ctx.QueryParam("email")
	if email == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email is required",
		})
	}

	var user models.User
	if err := database.CurrentDatabase.Where("email = ?", email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	err := c.authService.ResendConfirmation(email)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found or already verified",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "An error occurred while regenerating token",
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Nouveau email de confirmation envoyé",
	})
}
