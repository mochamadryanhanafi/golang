package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port         string
	DB           *gorm.DB
	Redis        *redis.Client
	JwtSecret    string
	SmtpHost     string
	SmtpPort     string
	SmtpUser     string
	SmtpPassword string
	AppEmail     string
}

var AppConfig *Config

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env variables")
	}

	// DB config
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect to database: %v", err)
	}

	// Redis config
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Load config values
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	AppConfig = &Config{
		Port:         port,
		DB:           db,
		Redis:        redisClient,
		JwtSecret:    os.Getenv("JWT_SECRET"),
		SmtpHost:     os.Getenv("MAIL_HOST"),
		SmtpPort:     os.Getenv("MAIL_PORT"),
		SmtpUser:     os.Getenv("MAIL_USERNAME"),
		SmtpPassword: os.Getenv("MAIL_PASSWORD"),
		AppEmail:     os.Getenv("MAIL_FROM"),
	}

	fmt.Println("✅ Configuration loaded successfully")
	return AppConfig
}
