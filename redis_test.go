package juice_cache

import (
	"context"
	"errors"
	cache2 "github.com/eatmoreapple/juice/cache"
	"github.com/redis/go-redis/v9"
	"testing"
)

var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func TestRedisCache(t *testing.T) {
	cache := NewRedisCache(client)
	if err := cache.Set(context.Background(), "test", 1); err != nil {
		t.Error(err)
	}
	var i int
	if err := cache.Get(context.Background(), "test", &i); err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Error("i != 1")
	}
	if err := cache.Flush(context.Background()); err != nil {
		t.Error(err)
	}
	if err := cache.Get(context.Background(), "test", &i); err != nil {
		if !errors.Is(err, cache2.ErrCacheNotFound) {
			t.Error(err)
		}
	}
}
