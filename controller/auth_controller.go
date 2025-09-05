package controller

import (
	"auth-service/model"
	"auth-service/service"
	"auth-service/utils"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthController(svc *service.AuthService, validate *validator.Validate) *AuthController {
	return &AuthController{
		authService: svc,
		validate:    validate,
	}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var input model.RegisterInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := ac.authService.Register(r.Context(), input)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			utils.WriteError(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered. Please check your email for the verification OTP.",
	})
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := ac.authService.Login(r.Context(), input)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			utils.WriteError(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, tokens)
}

func (ac *AuthController) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var input model.VerifyOTPInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := ac.authService.VerifyOTP(r.Context(), input.Email, input.OTP)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			utils.WriteError(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, tokens)
}

func (ac *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input model.RefreshTokenInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := ac.authService.RefreshToken(r.Context(), input.Token)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			utils.WriteError(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, tokens)
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	var input model.LogoutInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := ac.authService.Logout(r.Context(), input.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func (ac *AuthController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var input model.ForgotPasswordInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := ac.authService.ForgotPassword(r.Context(), input.Email)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			utils.WriteError(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "If a user with that email exists, a password reset link has been sent."})
}

func (ac *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input model.ResetPasswordInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := ac.authService.ResetPassword(r.Context(), input)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			utils.WriteError(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password has been reset successfully."})
}

func (ac *AuthController) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var input model.ResendOTPInput
	if err := utils.DecodeAndValidate(r, &input, ac.validate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := ac.authService.ResendOTP(r.Context(), input.Email)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			utils.WriteError(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "A new OTP has been sent to your email."})
}

func (ac *AuthController) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(model.ContextKey("userID")).(string)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user ID from context")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to your protected profile!",
		"userID":  userID,
	})
}
