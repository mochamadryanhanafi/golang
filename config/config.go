package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port                       string
	DB                         *gorm.DB
	Redis                      *redis.Client
	JwtSecret                  string `validate:"required"`
	SmtpHost                   string `validate:"required"`
	SmtpPort                   string `validate:"required"`
	SmtpUser                   string `validate:"required"`
	SmtpPassword               string `validate:"required"`
	AppEmail                   string `validate:"required,email"`
	AccessTokenDuration        time.Duration
	RefreshTokenDuration       time.Duration
	OTPDuration                time.Duration
	ResetPasswordTokenDuration time.Duration
}

func parseIntWithDefault(strVal string, defaultVal int) int {
	if val, err := strconv.Atoi(strVal); err == nil {
		return val
	}
	return defaultVal
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

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
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	accessTokenMin := parseIntWithDefault(os.Getenv("ACCESS_TOKEN_DURATION_MINUTES"), 15)
	refreshTokenHours := parseIntWithDefault(os.Getenv("REFRESH_TOKEN_DURATION_HOURS"), 168) // 7 days
	otpMin := parseIntWithDefault(os.Getenv("OTP_DURATION_MINUTES"), 5)
	resetTokenMin := parseIntWithDefault(os.Getenv("RESET_TOKEN_DURATION_MINUTES"), 15)

	cfg := &Config{
		Port:                       port,
		DB:                         db,
		Redis:                      redisClient,
		JwtSecret:                  os.Getenv("JWT_SECRET"),
		SmtpHost:                   os.Getenv("MAIL_HOST"),
		SmtpPort:                   os.Getenv("MAIL_PORT"),
		SmtpUser:                   os.Getenv("MAIL_USERNAME"),
		SmtpPassword:               os.Getenv("MAIL_PASSWORD"),
		AppEmail:                   os.Getenv("MAIL_FROM"),
		AccessTokenDuration:        time.Duration(accessTokenMin) * time.Minute,
		RefreshTokenDuration:       time.Duration(refreshTokenHours) * time.Hour,
		OTPDuration:                time.Duration(otpMin) * time.Minute,
		ResetPasswordTokenDuration: time.Duration(resetTokenMin) * time.Minute,
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation error: %w", err)
	}

	log.Println("✅ Configuration loaded successfully")
	return cfg, nil
}
