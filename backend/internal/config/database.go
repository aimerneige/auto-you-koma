package config

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/aimerneige/auto-you-koma/internal/models"
)

var DB *gorm.DB

func InitDB() error {
	dbPath := "./data/auto-you-koma.db"

	// Ensure data directory exists
	if err := os.MkdirAll("./data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Character{},
		&models.Project{},
		&models.Script{},
		&models.Storyboard{},
		&models.Panel{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	fmt.Println("Database initialized successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}