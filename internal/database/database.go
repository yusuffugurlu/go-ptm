package database

import (
	"os"

	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitializeDb() {
	var err error

	Db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_CONNECTION_URL")), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to database ", err)
	}

	if err := Db.AutoMigrate(
		&models.User{},
		&models.Balance{},
		&models.Transaction{},
		&models.AuditLog{}); err != nil {
		logger.Log.Fatal("Failed to migrate database", err)
	}

	logger.Log.Info("Database connected and migrated successfully!")
}