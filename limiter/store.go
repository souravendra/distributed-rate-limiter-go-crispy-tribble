package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store interface {
	Incr(key string) (int64, error)
	Expire(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
}

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr string) *RedisStore {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisStore{client: client, ctx: ctx}
}

func (r *RedisStore) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisStore) Expire(key string, ttl time.Duration) error {
	return r.client.Expire(r.ctx, key, ttl).Err()
}

func (r *RedisStore) GetTTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}
