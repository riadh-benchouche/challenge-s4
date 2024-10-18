package main

import (
	"backend/database"
	"backend/routers"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"os"
)

var appRouters = []routers.Router{
	&routers.MainRoute{},
}

func main() {
	fmt.Println("Starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	fmt.Printf("APP_MODE: %s\n", os.Getenv("ENVIRONMENT"))
	if os.Getenv("ENVIRONMENT") == "development" {
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set("Access-Control-Allow-Origin", "*")
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				return next(c)
			}
		})
	}

	// init database
	newDB, err := database.InitDB()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	defer newDB.CloseDB()

	// auto migrate database
	err = newDB.AutoMigrate()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	routers.LoadRoutes(e, appRouters...)

	addr := "0.0.0.0:" + os.Getenv("PORT")
	e.Logger.Fatal(e.Start(addr))
	fmt.Printf("Listening on %s\n", addr)
}
