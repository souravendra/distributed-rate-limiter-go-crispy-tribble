package limiter

import "sync"

// Functional Options Pattern

type Limiter struct {
	Strategy RateLimiter
}

type Option func(*Limiter)

func WithStrategy(s RateLimiter) Option {
	return func(l *Limiter) {
		l.Strategy = s
	}
}

func NewLimiter(opts ...Option) *Limiter {
	l := &Limiter{}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Singleton Limiter Instance
var (
	instance *Limiter
	once     sync.Once
)

func GetLimiter() *Limiter {
	once.Do(func() {
		store := NewRedisStore("localhost:6379")
		strategy := NewTokenBucket(store, 2, 1, 2) // 2 req/sec, burst 2
		instance = NewLimiter(WithStrategy(strategy))
	})
	return instance
}
