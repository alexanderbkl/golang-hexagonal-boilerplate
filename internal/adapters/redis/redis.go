package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/ports"
	"github.com/redis/go-redis/v9"
)

// RedisRepository implements the CacheRepository interface
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new Redis repository
func NewRedisRepository(client *redis.Client) ports.CacheRepository {
	return &RedisRepository{
		client: client,
	}
}

// Set stores a value in the cache
func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

// Get retrieves a value from the cache
func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a value from the cache
func (r *RedisRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in the cache
func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
