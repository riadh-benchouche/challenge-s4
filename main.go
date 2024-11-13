package main

import (
	"backend/database"
	"backend/faker"
	"backend/routers"
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
	&routers.AssociationRouter{},
	&routers.CategoryRouter{},
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
	defer database.CloseDB(newDB)

	err = newDB.AutoMigrate()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	routers.LoadRoutes(e, appRouters...)

	// Serve static files for Flutter web
	// e.Static("/app", utils.GetEnv("FLUTTER_BUILD_PATH", "flutter_build")+"/web")

	e.Static("/public", "public")

	faker.GenerateFakeData(newDB)

	addr := "0.0.0.0:3000"
	e.Logger.Fatal(e.Start(addr))
	fmt.Printf("Listening on %s\n", addr)
}
