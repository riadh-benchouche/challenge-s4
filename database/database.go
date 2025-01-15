package database

import (
	"backend/models"
	"fmt"
	"log"
	"os"
	"time"

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
		SSLMode:  "require",
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

func migrateEmailVerification(db *gorm.DB) error {
	type User struct {
		EmailVerifiedAt   *time.Time
		VerificationToken string
		TokenExpiresAt    *time.Time
	}

	// Vérifier si les colonnes existent déjà
	if !db.Migrator().HasColumn(&models.User{}, "email_verified_at") {
		err := db.Migrator().AddColumn(&models.User{}, "email_verified_at")
		if err != nil {
			return fmt.Errorf("failed to add email_verified_at column: %v", err)
		}
	}

	if !db.Migrator().HasColumn(&models.User{}, "verification_token") {
		err := db.Migrator().AddColumn(&models.User{}, "verification_token")
		if err != nil {
			return fmt.Errorf("failed to add verification_token column: %v", err)
		}
		// Ajouter l'index unique
		err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_verification_token ON users(verification_token) WHERE verification_token IS NOT NULL").Error
		if err != nil {
			return fmt.Errorf("failed to create verification_token index: %v", err)
		}
	}

	if !db.Migrator().HasColumn(&models.User{}, "token_expires_at") {
		err := db.Migrator().AddColumn(&models.User{}, "token_expires_at")
		if err != nil {
			return fmt.Errorf("failed to add token_expires_at column: %v", err)
		}
	}

	fmt.Println("✅ Email verification migration completed")
	return nil
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

func InitTestDB() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/backend_test"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err = db.AutoMigrate(Models...); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	CurrentDatabase = db
	return db, nil
}
