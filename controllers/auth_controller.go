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

		if errors.Is(err, coreErrors.ErrEmailNotVerified) {
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

	fmt.Printf("Token Firebase reçu dans la requête: %s\n", jsonBody.FirebaseToken)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(jsonBody)
	if err != nil {
		validationErrors := utils.GetValidationErrors(err.(validator.ValidationErrors), jsonBody)
		return ctx.JSON(http.StatusUnprocessableEntity, validationErrors)
	}

	result, err := c.authService.Register(jsonBody)
	if err != nil {
		if errors.Is(err, coreErrors.ErrEmailAlreadyExists) {
			return ctx.String(http.StatusConflict, "Email already used")
		}
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	// Générer le lien de confirmation avec le token
	confirmationLink := fmt.Sprintf("https://invooce.online/auth/confirm?token=%s", result.User.VerificationToken)
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

	if err := utils.SendEmail(result.User.Email, subject, body); err != nil {
		ctx.Logger().Error("Erreur lors de l'envoi de l'email :", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to send confirmation email"})
	}

	var newUser models.User
	if err := database.CurrentDatabase.Where("email = ?", jsonBody.Email).First(&newUser).Error; err != nil {
		fmt.Printf("Erreur lors de la récupération de l'utilisateur: %v\n", err)
	} else {
		fmt.Printf("Firebase Token dans la BDD: %s\n", newUser.FirebaseToken)

		if newUser.FirebaseToken != "" {
			fmt.Printf("Tentative d'envoi de notification avec le token: %s\n", newUser.FirebaseToken)
			notificationErr := utils.SendNotification(newUser.FirebaseToken, "Bienvenue sur VotreApp !",
				fmt.Sprintf("Un email de confirmation a été envoyé à %s. Vérifiez votre boîte mail.", newUser.Email))
			if notificationErr != nil {
				fmt.Printf("Erreur lors de l'envoi de la notification: %v\n", notificationErr)
				ctx.Logger().Error("Erreur lors de l'envoi de la notification push :", notificationErr)
			} else {
				fmt.Println("Notification push envoyée avec succès")
			}
		} else {
			fmt.Println("Token Firebase absent de la BDD")
		}
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"message": "Inscription réussie. Veuillez vérifier votre email pour confirmer votre compte",
	})
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
		switch {
		case errors.Is(err, coreErrors.ErrInvalidToken):
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

	// Envoyer une notification push à l'utilisateur
	if user.FirebaseToken != "" {
		notificationErr := utils.SendNotification(user.FirebaseToken, "Compte confirmé",
			"Votre compte a été confirmé avec succès. Vous pouvez maintenant vous connecter.")
		if notificationErr != nil {
			ctx.Logger().Error("Erreur lors de l'envoi de la notification push :", notificationErr)
		}
	}

	return ctx.HTML(http.StatusOK, `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<title>Email confirmé</title>
			<style>
				body { font-family: Arial, sans-serif; background-color: #d4edda; color: #007bff; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
				.container { text-align: center; padding: 20px; background-color: #ffffff; border: 1px solid #c3e6cb; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); }
				button { margin-top: 10px; padding: 10px 20px; background-color: #c3e6cb; color: #007bff; border: none; border-radius: 5px; cursor: pointer; }
				button:hover { background-color: #d4edda; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Félicitations</h1>
				<p>Votre email a été confirmé avec succès ! Vous pouvez maintenant vous connecter à votre compte.</p>
			</div>
		</body>
		</html>
	`)
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
		switch {
		case errors.Is(err, coreErrors.ErrUserNotFound):
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

func (c *AuthController) RefreshToken(ctx echo.Context) error {
	// Structure pour parser le JSON
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Décoder le JSON du body
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	fmt.Println("Refresh token reçu:", request.RefreshToken) // Pour debug

	if request.RefreshToken == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Refresh token is required",
		})
	}

	tokens, err := c.authService.RefreshToken(request.RefreshToken)
	if err != nil {
		switch err {
		case coreErrors.ErrInvalidToken:
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired refresh token",
			})
		case coreErrors.ErrUserNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Internal server error",
			})
		}
	}

	return ctx.JSON(http.StatusOK, tokens)
}

func (c *AuthController) ForgotPassword(ctx echo.Context) error {
	// Analyse des données JSON de la requête
	var request struct {
		Email string `json:"email"`
	}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if request.Email == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email is required",
		})
	}

	// Appeler le service pour générer le token de réinitialisation
	resetToken, err := c.authService.GeneratePasswordResetToken(request.Email)
	if err != nil {
		switch {
		case errors.Is(err, coreErrors.ErrUserNotFound):
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		default:
			ctx.Logger().Error(err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "An error occurred while generating reset token",
			})
		}
	}

	// Créer un lien de réinitialisation
	resetLink := fmt.Sprintf("https://invooce.online/auth/reset-password?token=%s", resetToken)
	subject := "Réinitialisation de votre mot de passe"
	body := fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="fr">
    <head>
        <meta charset="UTF-8">
        <title>Réinitialisation de mot de passe</title>
    </head>
    <body>
        <h2>Demande de réinitialisation de mot de passe</h2>
        <p>Vous avez demandé la réinitialisation de votre mot de passe. Cliquez sur le lien ci-dessous pour définir un nouveau mot de passe :</p>
        <a href="%s">Réinitialiser mon mot de passe</a>
        <p>Si vous n'avez pas demandé cette réinitialisation, vous pouvez ignorer cet email.</p>
    </body>
    </html>
    `, resetLink)

	// Envoyer un email
	if err := utils.SendEmail(request.Email, subject, body); err != nil {
		ctx.Logger().Error("Erreur lors de l'envoi de l'email :", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to send reset email"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Email de réinitialisation envoyé avec succès",
	})
}

func (c *AuthController) ResetPassword(ctx echo.Context) error {
	// Récupérer les paramètres du formulaire
	token := ctx.FormValue("token")
	newPassword := ctx.FormValue("new_password")

	// Vérifier si le token et le mot de passe sont présents
	if token == "" || newPassword == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Token and new password are required",
		})
	}

	// Passer ces valeurs à votre service pour réinitialiser le mot de passe
	err := c.authService.ResetPassword(token, newPassword)
	if err != nil {
		switch {
		case errors.Is(err, coreErrors.ErrInvalidToken):
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired token",
			})
		case errors.Is(err, coreErrors.ErrUserNotFound):
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		default:
			ctx.Logger().Error(err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "An error occurred while resetting password",
			})
		}
	}

	return c.ResetPasswordSuccess(ctx)

}

func (c *AuthController) ResetPasswordForm(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	if token == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Token is required",
		})
	}

	return ctx.HTML(http.StatusOK, fmt.Sprintf(`
        <html>
            <body>
				<div class="container">
					<h1>Réinitialisation du mot de passe</h1>
					<form method="POST" action="/auth/reset-password">
						<input type="hidden" name="token" value="%s" />
						<label for="new_password">Nouveau mot de passe :</label>
						<input type="password" name="new_password" id="new_password" required />
						<button type="submit">Réinitialiser</button>
					</form>
				</div>
            </body>
			          <style>
                body {
                    font-family: 'Arial', sans-serif;
                    background-color: #f2f2f2;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    border-radius: 8px;
                    padding: 20px;
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
                    width: 100%;
                    max-width: 400px;
                }
                h1 {
                    text-align: center;
                    font-size: 24px;
                    margin-bottom: 20px;
                }
                input {
                    width: 100%;
                    padding: 10px;
                    margin: 10px 0;
                    border-radius: 4px;
                    border: 1px solid #ddd;
                    box-sizing: border-box;
                }
                button {
                    width: 100%;
                    padding: 10px;
                    background-color: #007bff;
                    border: none;
                    border-radius: 4px;
                    color: white;
                    font-size: 16px;
                    cursor: pointer;
                }
                button:hover {
                    background-color: #0056b3;
                }
                .message {
                    text-align: center;
                    font-size: 16px;
                    margin-top: 20px;
                }
            </style>
        </html>
    `, token))
}

func (c *AuthController) ResetPasswordSuccess(ctx echo.Context) error {
	return ctx.HTML(http.StatusOK, `
        <!DOCTYPE html>
        <html lang="fr">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Réinitialisation du Mot de Passe - Succès</title>
            <style>
                body {
                    font-family: 'Arial', sans-serif;
                    background-color: #f2f2f2;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    border-radius: 8px;
                    padding: 20px;
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
                    width: 100%;
                    max-width: 400px;
                }
                h1 {
                    text-align: center;
                    font-size: 24px;
                    margin-bottom: 20px;
                }
                p {
                    text-align: center;
                    font-size: 18px;
                    color: #28a745;
                    margin-top: 20px;
                }
                .button-container {
                    text-align: center;
                    margin-top: 30px;
                }
                .back-button {
                    padding: 10px 20px;
                    background-color: #007bff;
                    border: none;
                    border-radius: 4px;
                    color: white;
                    font-size: 16px;
                    cursor: pointer;
                    text-decoration: none;
                }
                .back-button:hover {
                    background-color: #0056b3;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>Mot de Passe Réinitialisé</h1>
                <p>Votre mot de passe a été réinitialisé avec succès !</p>
            </div>
        </body>
        </html>
    `)
}
