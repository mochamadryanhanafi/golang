package controller

import (
	"auth-service/model"
	"auth-service/service"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(svc *service.AuthService) *AuthController {
	return &AuthController{AuthService: svc}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var input model.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := ac.AuthService.Register(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered. Please verify OTP sent to your email."))
}

func (ac *AuthController) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var input model.VerifyOTPInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := ac.AuthService.VerifyOTP(r.Context(), input.Email, input.OTP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := ac.AuthService.Login(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write([]byte("OTP sent to your email."))
}
