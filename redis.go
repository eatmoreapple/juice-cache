package juice_cache

import (
	"context"
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
		return r.client.HSet(ctx, r.uuid, key, value).Err()
	}
}

// Get gets the value for the key.
func (r *RedisCache) Get(ctx context.Context, key string, dst any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// dst must implement encoding.BinaryUnmarshaler if it is not basic type.
		err := r.client.HGet(ctx, r.uuid, key).Scan(dst)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return cache.ErrCacheNotFound
			}
		}
		return err
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
