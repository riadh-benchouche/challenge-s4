package main

import (
	"backend/controllers"
	"backend/database"
	"backend/routers"
	"backend/services"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var appRouters = []routers.Router{
	&routers.HelloRouter{},
	&routers.UserRouter{},
	&routers.AuthRouter{},
}

func main() {
	fmt.Println("Starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	// Middleware CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	fmt.Printf("APP_MODE: %s\n", os.Getenv("ENVIRONMENT"))

	// Initialisation de la base de données
	newDB, err := database.InitDB()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	routers.LoadRoutes(e, appRouters...)

	// Chargement des routes avec AssociationController
	associationService := services.NewAssociationService(newDB) // Service pour les associations
	associationController := controllers.NewAssociationController(associationService)
	routers.SetupAssociationRoutes(e, associationController) // Charger les routes de l'association

	// Démarrer le serveur
	addr := "0.0.0.0:3000" // Port fixe pour éviter d'utiliser une variable d'environnement ici
	e.Logger.Fatal(e.Start(addr))
	fmt.Printf("Listening on %s\n", addr)
}
