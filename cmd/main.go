package main

import (
	"auth-service/config"
	"auth-service/routes"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()
	config.InitDB()
	r := chi.NewRouter()
	routes.SetupRoutes(r)

	log.Printf("Server started on :%s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
