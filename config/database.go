// di file: config/database.go

package config

import (
	"fmt"
	"log"
	"ticketing/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Membuka koneksi ke database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")

	// AutoMigrate models
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db, nil
}

// AutoMigrate menjalankan migrasi otomatis untuk model-model yang telah ditentukan
func AutoMigrate(db *gorm.DB) error {
	// Migrasi semua model yang digunakan
	return db.AutoMigrate(&model.User{}, &model.Event{}, &model.Ticket{})
}
