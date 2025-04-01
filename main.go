package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// --- Strategy Pattern ---
type RateLimiter interface {
	Allow(key string) bool
}

// Token Bucket implementation using Redis

// --- Adapter Pattern for Store ---
type Store interface {
	Incr(key string) (int64, error)
	Expire(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
}

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
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

// Token Bucket Strategy

type TokenBucket struct {
	store     Store
	rate      int           // tokens per interval
	interval  time.Duration // e.g., 1 second
	burst     int
	keyPrefix string
}

func NewTokenBucket(store Store, rate int, interval time.Duration, burst int) *TokenBucket {
	return &TokenBucket{
		store:     store,
		rate:      rate,
		interval:  interval,
		burst:     burst,
		keyPrefix: "rate:limiter:",
	}
}

func (tb *TokenBucket) Allow(key string) bool {
	redisKey := tb.keyPrefix + key
	count, err := tb.store.Incr(redisKey)
	if err != nil {
		log.Println("Redis INCR error:", err)
		return true // fail open
	}

	if count == 1 {
		// Set TTL
		err := tb.store.Expire(redisKey, tb.interval)
		if err != nil {
			log.Println("Redis EXPIRE error:", err)
		}
		ttl, _ := tb.store.GetTTL(redisKey)
		fmt.Printf("New TTL set: %v\n", ttl)
	}

	fmt.Printf("Key: %s, Count: %d\n", redisKey, count)

	return int(count) <= tb.burst
}

// --- Functional Options Pattern ---
type Limiter struct {
	strategy RateLimiter
}

type Option func(*Limiter)

func WithStrategy(s RateLimiter) Option {
	return func(l *Limiter) {
		l.strategy = s
	}
}

func NewLimiter(opts ...Option) *Limiter {
	l := &Limiter{}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// --- Singleton Pattern ---
var (
	limiterInstance *Limiter
	once            sync.Once
)

func GetLimiter() *Limiter {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		ctx := context.Background()
		store := &RedisStore{client: rdb, ctx: ctx}
		strategy := NewTokenBucket(store, 2, time.Second, 2)
		limiterInstance = NewLimiter(WithStrategy(strategy))
	})
	return limiterInstance
}

// --- Middleware Pattern ---
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := "test-client" // r.RemoteAddr
		if !GetLimiter().strategy.Allow(key) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- Main server ---
func main() {
	http.Handle("/", RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request allowed: %s\n", time.Now().Format(time.RFC3339))
	})))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
