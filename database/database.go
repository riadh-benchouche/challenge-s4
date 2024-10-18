package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
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

type DB struct {
	DB     *gorm.DB
	Config Config
}

func (db *DB) Connect() error {
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Config.Host, db.Config.Port, db.Config.User, db.Config.Password, db.Config.Name, db.Config.SSLMode)

	var err error
	db.DB, err = gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	sqlDB, _ := db.DB.DB()
	return sqlDB.Close()
}

func (db *DB) CloseDB() {
	err := db.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func InitDB() (*DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var DefaultConfig = Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DATABASE"),
		Port:     os.Getenv("POSTGRES_PORT"),
		SSLMode:  "disable",
	}

	newDB := DB{Config: DefaultConfig}

	errorConnection := newDB.Connect()
	if errorConnection != nil {
		return nil, err
	}

	CurrentDatabase = newDB.DB
	return &newDB, nil
}
