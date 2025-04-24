package db

import (
	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dbConnectionURL string) {
	db, err := gorm.Open(postgres.Open(dbConnectionURL), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", err)
	}

	DB = db
	logger.Log.Info("Database connection established successfully")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Balance{},
		&models.Transaction{},); err != nil {
		logger.Log.Fatal("Failed to migrate database", err)
	}

	logger.Log.Info("Database migrated successfully")
}