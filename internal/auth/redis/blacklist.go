package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenBlacklist struct {
	client *redis.Client
	prefix string
}

func NewRedisTokenBlacklist(addr, password string, db int) *RedisTokenBlacklist {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisTokenBlacklist{
		client: client,
		prefix: "blacklist:",
	}
}

// Add adiciona o token à blacklist com expiração automática no Redis
func (b *RedisTokenBlacklist) Add(ctx context.Context, token string, duration time.Duration) error {
	key := b.prefix + token
	return b.client.Set(ctx, key, "revoked", duration).Err()
}

// IsBlacklisted verifica se o token está na blacklist
func (b *RedisTokenBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := b.prefix + token
	result, err := b.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
