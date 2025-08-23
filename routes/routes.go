package routes

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/middleware"
	"auth-service/repository"
	"auth-service/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	// Ambil koneksi database dari config
	db := config.GetDB()

	// Inisialisasi repositori dan service
	userRepo := repository.NewUserRepo(db)
	redisRepo := repository.NewRedisRepo() // asumsi tidak butuh DB
	authService := service.NewAuthService(userRepo, redisRepo)
	authController := controller.NewAuthController(authService)

	// Public routes
	r.Post("/register", authController.Register)
	r.Post("/login", authController.Login)
	r.Post("/verify-otp", authController.VerifyOTP)

	// Protected route
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTMiddleware)
		protected.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You are authenticated!"))
		})
	})
}
