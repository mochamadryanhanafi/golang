package service

import (
	"auth-service/config"
	"auth-service/model"
	"auth-service/repository"
	"auth-service/utils"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  *repository.UserRepo
	RedisRepo *repository.RedisRepo
}

func NewAuthService(userRepo *repository.UserRepo, redisRepo *repository.RedisRepo) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		RedisRepo: redisRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, input model.RegisterInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		IsVerified:   false,
	}

	err = s.UserRepo.Create(ctx, &user)
	if err != nil {
		return err
	}

	otp := utils.GenerateOTP()
	err = s.RedisRepo.SaveOTP(ctx, user.Email, otp, 5*time.Minute)
	if err != nil {
		return err
	}

	return utils.SendOTPEmail(user.Email, otp)
}

func (s *AuthService) Login(ctx context.Context, input model.LoginInput) error {
	user, err := s.UserRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return errors.New("invalid credentials")
	}

	otp := utils.GenerateOTP()
	err = s.RedisRepo.SaveOTP(ctx, user.Email, otp, 5*time.Minute)
	if err != nil {
		return err
	}

	return utils.SendOTPEmail(user.Email, otp)
}

func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) (map[string]string, error) {
	isValid := s.RedisRepo.VerifyOTP(ctx, email, otp)
	if !isValid {
		return nil, errors.New("invalid or expired OTP")
	}

	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.IsVerified {
		user.IsVerified = true
		s.UserRepo.Update(ctx, user)
	}

	accessToken, err := utils.GenerateJWT(user.ID, config.AppConfig.JwtSecret, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = s.RedisRepo.SaveRefreshToken(ctx, user.ID.String(), refreshToken, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}
