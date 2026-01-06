package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	CloudinaryName      string
	CloudinaryAPIKey    string
	CLoudinaryAPISecret string
	ServerPort          string
	Environment         string
	GinMode             string
	Domain              string
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	EncryptionKey       string
}

var Config *AppConfig

func LoadConfig() {

	if os.Getenv("ENVIRONMENT") != "production" {
		godotenv.Load()
	}

	Config = &AppConfig{
		DBHost:              os.Getenv("DATABASE_HOST"),
		DBPort:              os.Getenv("DATABASE_PORT"),
		DBUser:              os.Getenv("DATABASE_USERNAME"),
		DBPassword:          os.Getenv("DATABASE_PASSWORD"),
		DBName:              os.Getenv("DATABASE_NAME"),
		CloudinaryName:      os.Getenv("CLOUDINARY_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CLoudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		ServerPort:          os.Getenv("PORT"),
		Environment:         os.Getenv("ENVIRONMENT"),
		GinMode:             os.Getenv("GIN_MODE"),
		Domain:              os.Getenv("DOMAIN"),
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:   os.Getenv("GOOGLE_REDIRECT_URL"),
		EncryptionKey:       os.Getenv("ENCRYPTION_KEY"),
	}

}
