package repository

import (
	"auth-service/config"
	"context"
	"fmt"
	"time"
)

type RedisRepo struct{}

func NewRedisRepo() *RedisRepo {
	return &RedisRepo{}
}

func (r *RedisRepo) SaveOTP(ctx context.Context, email, otp string, ttl time.Duration) error {
	key := fmt.Sprintf("otp:%s", email)
	return config.AppConfig.Redis.Set(ctx, key, otp, ttl).Err()
}

func (r *RedisRepo) VerifyOTP(ctx context.Context, email, otp string) bool {
	key := fmt.Sprintf("otp:%s", email)
	val, err := config.AppConfig.Redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	return val == otp
}

func (r *RedisRepo) SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%s", userID)
	return config.AppConfig.Redis.Set(ctx, key, token, ttl).Err()
}

func (r *RedisRepo) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("refresh:%s", userID)
	return config.AppConfig.Redis.Get(ctx, key).Result()
}

func (r *RedisRepo) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf("refresh:%s", userID)
	return config.AppConfig.Redis.Del(ctx, key).Err()
}
