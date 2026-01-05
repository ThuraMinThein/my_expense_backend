package db

import (
	"fmt"

	"github.com/ThuraMinThein/my_expense_backend/config"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit(migrateDatabase bool) error {
	config := config.Config
	sslMode := "require"
	if config.Environment == "development" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
		sslMode,
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if migrateDatabase {
		DB.AutoMigrate(&models.User{}, &models.UserToken{})
	}

	return nil
}
