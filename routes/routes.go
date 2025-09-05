package routes

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux, authController *controller.AuthController, cfg *config.Config) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Post("/verify-otp", authController.VerifyOTP)
		r.Post("/token/refresh", authController.RefreshToken)
		r.Post("/forgot-password", authController.ForgotPassword)
		r.Post("/reset-password", authController.ResetPassword)
		r.Post("/resend-otp", authController.ResendOTP)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(cfg.JwtSecret))

		r.Post("/auth/logout", authController.Logout)
		r.Get("/profile", authController.GetProfile)
	})
}
