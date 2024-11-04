package database

import (
	"backend/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var CurrentDatabase *gorm.DB

type Filter struct {
	Value    interface{} `json:"value" validate:"required"`
	Operator string      `json:"operator" validate:"required,oneof= != > < >= <= ="`
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

var Models = []interface{}{
	&models.User{},
	&models.Association{},
	&models.Membership{},
	&models.Message{},
	&models.Category{},
	&models.Event{},
	&models.Participation{},
}

// InitDB initialise la base de données et effectue la migration
func InitDB() (*gorm.DB, error) {
	fmt.Println("🚀 Initializing database...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Configuration de la base de données
	config := Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DATABASE"),
		Port:     os.Getenv("POSTGRES_PORT"),
		SSLMode:  "disable",
	}

	// Chaîne de connexion
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)

	// Connexion à la base de données
	fmt.Println("⏳ Waiting for database connection...")
	db, err := gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Vérifier la connexion
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("🎉 Database connected!")

	// Migrer les modèles
	err = db.AutoMigrate(Models...)
	if err != nil {
		return nil, err
	}

	// Stocker la base de données actuelle dans CurrentDatabase
	CurrentDatabase = db
	return db, nil
}

// CloseDB ferme la connexion à la base de données
func CloseDB(db *gorm.DB) {
	fmt.Println("🚨 Closing database connection...")
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Error obtaining DB:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Println("Error closing DB:", err)
	}
}
