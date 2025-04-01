package limiter

import (
	"fmt"
	"log"
	"time"
)

type RateLimiter interface {
	Allow(key string) bool
}

type TokenBucket struct {
	store     Store
	rate      int
	interval  time.Duration
	burst     int
	keyPrefix string
}

func NewTokenBucket(store Store, rate int, intervalSeconds int, burst int) *TokenBucket {
	return &TokenBucket{
		store:     store,
		rate:      rate,
		interval:  time.Duration(intervalSeconds) * time.Second,
		burst:     burst,
		keyPrefix: "rate:limiter:",
	}
}

func (tb *TokenBucket) Allow(key string) bool {
	redisKey := tb.keyPrefix + key
	count, err := tb.store.Incr(redisKey)
	if err != nil {
		log.Println("Redis INCR error:", err)
		return true
	}

	if count == 1 {
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
