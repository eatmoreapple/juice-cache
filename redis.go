package juice_cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eatmoreapple/juice/cache"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisCache is a redis cache implementation.
// It is used redis hash to store cache data.
// Each cache instance must have a unique uuid.
type RedisCache struct {
	client *redis.Client
	uuid   string
}

// Set sets the value for the key.
func (r *RedisCache) Set(ctx context.Context, key string, value any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return r.client.HSet(ctx, r.uuid, key, data).Err()
	}
}

// Get gets the value for the key.
func (r *RedisCache) Get(ctx context.Context, key string, dst any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var data []byte
		err := r.client.HGet(ctx, r.uuid, key).Scan(&data)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return cache.ErrCacheNotFound
			}
			return err
		}
		return json.Unmarshal(data, dst)
	}
}

// Flush flushes all the cache.
func (r *RedisCache) Flush(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return r.client.Del(ctx, r.uuid).Err()
	}
}

// NewRedisCache returns a new redis cache instance.
func NewRedisCache(client *redis.Client) cache.Cache {
	return &RedisCache{client: client, uuid: uuid.NewString()}
}
