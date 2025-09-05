package main

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/repository"
	"auth-service/routes"
	"auth-service/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	validate := validator.New()

	userRepo := repository.NewUserRepo(cfg.DB)
	redisRepo := repository.NewRedisRepo(cfg.Redis)

	authService := service.NewAuthService(userRepo, redisRepo, cfg)

	authController := controller.NewAuthController(authService, validate)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.SetupRoutes(r, authController, cfg)

	log.Printf("🚀 Server started on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
