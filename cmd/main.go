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
	// Log-log ini akan membantu Anda melacak titik kegagalan.
	// Perhatikan output di terminal Anda saat menjalankan program.
	log.Println("--- [Step 1] Memuat konfigurasi ---")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Tidak dapat memuat konfigurasi: %v", err)
	}
	log.Println("--- [Step 1] Konfigurasi berhasil dimuat ---")

	log.Println("--- [Step 2] Menginisialisasi validator ---")
	validate := validator.New()
	log.Println("--- [Step 2] Validator berhasil diinisialisasi ---")

	log.Println("--- [Step 3] Menginisialisasi repositories ---")
	userRepo := repository.NewUserRepo(cfg.DB)
	redisRepo := repository.NewRedisRepo(cfg.Redis)
	log.Println("--- [Step 3] Repositories berhasil diinisialisasi ---")

	log.Println("--- [Step 4] Menginisialisasi services ---")
	authService := service.NewAuthService(userRepo, redisRepo, cfg)
	log.Println("--- [Step 4] Services berhasil diinisialisasi ---")

	log.Println("--- [Step 5] Menginisialisasi controllers ---")
	authController := controller.NewAuthController(authService, validate)
	log.Println("--- [Step 5] Controllers berhasil diinisialisasi ---")

	log.Println("--- [Step 6] Menyiapkan router dan middleware ---")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	log.Println("--- [Step 6] Router dan middleware berhasil disiapkan ---")

	log.Println("--- [Step 7] Menyiapkan rute ---")
	routes.SetupRoutes(r, authController, cfg)
	log.Println("--- [Step 7] Rute berhasil disiapkan ---")

	log.Println("--- [Step 8] Memulai server ---")
	log.Printf("🚀 Server dimulai pada port :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("FATAL: Gagal memulai server: %v", err)
	}
}
