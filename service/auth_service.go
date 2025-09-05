package service

import (
	"auth-service/config"
	"auth-service/model"
	"auth-service/repository"
	"auth-service/utils"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepo
	redisRepo *repository.RedisRepo
	cfg       *config.Config
}

func NewAuthService(userRepo *repository.UserRepo, redisRepo *repository.RedisRepo, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		redisRepo: redisRepo,
		cfg:       cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, input model.RegisterInput) error {
	_, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil {
		return model.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	user := model.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		IsVerified:   false,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	return s.sendOTP(ctx, user.Email)
}

func (s *AuthService) Login(ctx context.Context, input model.LoginInput) (map[string]string, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, model.ErrInvalidCredentials
	}

	if !user.IsVerified {
		s.sendOTP(ctx, user.Email)
		return nil, model.ErrAccountNotVerified
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) (map[string]string, error) {
	isValid, err := s.redisRepo.VerifyOTP(ctx, email, otp)
	if err != nil {
		return nil, fmt.Errorf("could not verify otp from redis: %w", err)
	}
	if !isValid {
		return nil, model.ErrInvalidOTP
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	if !user.IsVerified {
		user.IsVerified = true
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("could not update user verification status: %w", err)
		}
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error) {
	userID, err := s.redisRepo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	if err := s.redisRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		fmt.Printf("warning: could not delete old refresh token: %v", err)
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.redisRepo.DeleteRefreshToken(ctx, refreshToken)
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil
	}

	token := utils.GenerateSecureRandomString(32)
	if err := s.redisRepo.SaveResetToken(ctx, token, user.Email, s.cfg.ResetPasswordTokenDuration); err != nil {
		return fmt.Errorf("could not save reset token: %w", err)
	}

	return utils.SendResetPasswordEmail(user.Email, token, s.cfg)
}

func (s *AuthService) ResetPassword(ctx context.Context, input model.ResetPasswordInput) error {
	email, err := s.redisRepo.GetEmailByResetToken(ctx, input.Token)
	if err != nil {
		return model.ErrInvalidToken
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return model.ErrUserNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash new password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("could not update password: %w", err)
	}

	s.redisRepo.DeleteResetToken(ctx, input.Token)
	return nil
}

func (s *AuthService) ResendOTP(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return model.ErrUserNotFound
	}
	if user.IsVerified {
		return model.NewAppError(400, "Account is already verified")
	}
	return s.sendOTP(ctx, email)
}

// --- Helper Functions ---

func (s *AuthService) sendOTP(ctx context.Context, email string) error {
	otp := utils.GenerateOTP(6)
	if err := s.redisRepo.SaveOTP(ctx, email, otp, s.cfg.OTPDuration); err != nil {
		return fmt.Errorf("could not save OTP: %w", err)
	}
	return utils.SendOTPEmail(email, otp, s.cfg)
}

func (s *AuthService) generateTokens(ctx context.Context, user *model.User) (map[string]string, error) {
	accessToken, err := utils.GenerateJWT(user.ID.String(), s.cfg.JwtSecret, s.cfg.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("could not generate access token: %w", err)
	}

	refreshToken := uuid.New().String()
	if err := s.redisRepo.SaveRefreshToken(ctx, refreshToken, user.ID.String(), s.cfg.RefreshTokenDuration); err != nil {
		return nil, fmt.Errorf("could not save refresh token: %w", err)
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}
