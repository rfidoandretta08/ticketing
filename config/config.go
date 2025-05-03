package config

import (
	"fmt"
	"log"
	"os"

	"ticketing/model" // Sesuaikan dengan path model Anda

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() *Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}

func SetupDatabase(config *Config) *gorm.DB {
	// Menyiapkan koneksi ke database MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	log.Println("Successfully connected to the database")

	// Melakukan AutoMigrate untuk model yang ada
	err = db.AutoMigrate(
		&model.Event{},
		&model.Ticket{},
		&model.User{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrated successfully")
	return db
}
