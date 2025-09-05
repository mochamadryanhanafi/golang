package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) SaveOTP(ctx context.Context, email, otp string, ttl time.Duration) error {
	key := fmt.Sprintf("otp:%s", email)
	return r.client.Set(ctx, key, otp, ttl).Err()
}

func (r *RedisRepo) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	key := fmt.Sprintf("otp:%s", email)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // OTP tidak ditemukan (salah atau kedaluwarsa)
	} else if err != nil {
		return false, err // Error Redis lainnya
	}

	if val == otp {
		// Hapus OTP setelah berhasil diverifikasi
		r.client.Del(ctx, key)
		return true, nil
	}

	return false, nil
}

func (r *RedisRepo) SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%s", token) // Kunci berdasarkan token itu sendiri
	return r.client.Set(ctx, key, userID, ttl).Err()
}

// GetUserIDByRefreshToken mengambil User ID yang terkait dengan refresh token.
func (r *RedisRepo) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("refresh:%s", token)
	return r.client.Get(ctx, key).Result()
}

func (r *RedisRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh:%s", token)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisRepo) SaveResetToken(ctx context.Context, token, email string, ttl time.Duration) error {
	key := fmt.Sprintf("reset:%s", token)
	return r.client.Set(ctx, key, email, ttl).Err()
}

func (r *RedisRepo) GetEmailByResetToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("reset:%s", token)
	email, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	// Catatan: Token ini tidak dihapus di sini.
	// Service akan memanggil DeleteResetToken setelah password berhasil diubah.
	return email, nil
}

// DeleteResetToken menghapus token reset password dari Redis.
func (r *RedisRepo) DeleteResetToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("reset:%s", token)
	return r.client.Del(ctx, key).Err()
}
