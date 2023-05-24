package juice_cache

import (
	"context"
	"encoding"
	"errors"
	"github.com/eatmoreapple/juice/cache"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
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
		// if the value can be scanned by redis, then set it directly.
		if !redisBinaryMarshalerAble(value) {
			// otherwise, marshal it to json.
			value = jsonMarshalBinaryWrap(value)
		}
		return r.client.HSet(ctx, r.uuid, key, value).Err()
	}
}

// Get gets the value for the key.
func (r *RedisCache) Get(ctx context.Context, key string, dst any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var err error
		// if the value can be scanned by redis, then get it directly.
		if !redisBinaryUnmarshalerAble(dst) {
			// otherwise, unmarshal it from json.
			dst = jsonUnmarshalBinaryWrap(dst)
		}
		err = r.client.HGet(ctx, r.uuid, key).Scan(dst)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return cache.ErrCacheNotFound
			}
			return err
		}
		return nil
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

// redisBinaryUnmarshalerAble returns true if the value can be scanned by redis.
// see https://pkg.go.dev/github.com/redis/go-redis/v9#example-Client.Scan
func redisBinaryUnmarshalerAble(v any) bool {
	switch v.(type) {
	case *string, *[]byte, *int, *int8, *int16, *int32, *int64, *uint, *uint8,
		*uint16, *uint32, *uint64, *float32, *float64, *bool, *time.Time, *time.Duration,
		encoding.BinaryUnmarshaler, *net.IP:
		return true
	default:
		return false
	}
}

// redisBinaryMarshalerAble returns true if the value can be scanned by redis.
// see https://pkg.go.dev/github.com/redis/go-redis/v9#example-Client.Scan
func redisBinaryMarshalerAble(v any) bool {
	switch v.(type) {
	case *string, *[]byte, *int, *int8, *int16, *int32, *int64, *uint, *uint8,
		*uint16, *uint32, *uint64, *float32, *float64, *bool, *time.Time, *time.Duration,
		encoding.BinaryMarshaler, *net.IP:
		return true
	default:
		return false
	}
}

// NewRedisCache returns a new redis cache instance.
func NewRedisCache(client *redis.Client) cache.Cache {
	return &RedisCache{client: client, uuid: uuid.NewString()}
}
