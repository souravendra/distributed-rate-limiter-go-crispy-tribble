package limiter

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/souravendra/distributed-rate-limiter-go-crispy-tribble/config"
)

// --- Strategy Interface ---
type RateLimiter interface {
	Allow(key string) bool
}

var cfg = config.Get()

// --- Token Bucket ---
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
		keyPrefix: "rate:limiter:token:",
	}
}

func NewTokenBucketFromConfig(store Store) *TokenBucket {
	return NewTokenBucket(store, cfg.Rate, cfg.Interval, cfg.Burst)
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
		fmt.Printf("[Token] New TTL set: %v\n", ttl)
	}

	fmt.Printf("[Token] Key: %s, Count: %d\n", redisKey, count)
	return int(count) <= tb.burst
}

// --- Fixed Window ---
type FixedWindow struct {
	store     Store
	rate      int
	interval  time.Duration
	keyPrefix string
}

func NewFixedWindow(store Store, rate int, intervalSeconds int) *FixedWindow {
	return &FixedWindow{
		store:     store,
		rate:      rate,
		interval:  time.Duration(intervalSeconds) * time.Second,
		keyPrefix: "rate:limiter:fixed:",
	}
}

func NewFixedWindowFromConfig(store Store) *FixedWindow {
	return NewFixedWindow(store, cfg.Rate, cfg.Interval)
}

func (fw *FixedWindow) Allow(key string) bool {
	redisKey := fmt.Sprintf("%s%s:%d", fw.keyPrefix, key, time.Now().Unix()/int64(fw.interval.Seconds()))
	count, err := fw.store.Incr(redisKey)
	if err != nil {
		log.Println("Redis INCR error:", err)
		return true
	}
	if count == 1 {
		_ = fw.store.Expire(redisKey, fw.interval)
	}
	fmt.Printf("[Fixed] Key: %s, Count: %d\n", redisKey, count)
	return int(count) <= fw.rate
}

// --- Sliding Window Log ---
type SlidingWindowLog struct {
	store     Store
	rate      int
	interval  time.Duration
	keyPrefix string
	mu        sync.Mutex
	logs      map[string][]int64
}

func NewSlidingWindowLog(store Store, rate int, intervalSeconds int) *SlidingWindowLog {
	return &SlidingWindowLog{
		store:     store,
		rate:      rate,
		interval:  time.Duration(intervalSeconds) * time.Second,
		keyPrefix: "rate:limiter:log:",
		logs:      make(map[string][]int64),
	}
}

func NewSlidingWindowLogFromConfig(store Store) *SlidingWindowLog {
	return NewSlidingWindowLog(store, cfg.Rate, cfg.Interval)
}

func (sw *SlidingWindowLog) Allow(key string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now().Unix()
	windowStart := now - int64(sw.interval.Seconds())
	logKey := sw.keyPrefix + key

	entries := sw.logs[logKey]
	validEntries := []int64{}
	for _, ts := range entries {
		if ts > windowStart {
			validEntries = append(validEntries, ts)
		}
	}

	sw.logs[logKey] = append(validEntries, now)
	fmt.Printf("[SlidingLog] Key: %s, Count: %d\n", logKey, len(sw.logs[logKey]))
	return len(sw.logs[logKey]) <= sw.rate
}

// --- Leaky Bucket ---
type LeakyBucket struct {
	store     Store
	rate      int // leak rate per second
	interval  time.Duration
	burst     int
	keyPrefix string
}

func NewLeakyBucket(store Store, rate int, intervalSeconds int, burst int) *LeakyBucket {
	return &LeakyBucket{
		store:     store,
		rate:      rate,
		interval:  time.Duration(intervalSeconds) * time.Second,
		burst:     burst,
		keyPrefix: "rate:limiter:leaky:",
	}
}

func NewLeakyBucketFromConfig(store Store) *LeakyBucket {
	return NewLeakyBucket(store, cfg.Rate, cfg.Interval, cfg.Burst)
}

func (lb *LeakyBucket) Allow(key string) bool {
	redisKey := lb.keyPrefix + key
	count, err := lb.store.Incr(redisKey)
	if err != nil {
		log.Println("Redis INCR error:", err)
		return true
	}

	if count == 1 {
		_ = lb.store.Expire(redisKey, lb.interval)
	}

	fmt.Printf("[Leaky] Key: %s, Count: %d\n", redisKey, count)
	return int(count) <= lb.burst
}

// --- Moving Window Counter (Hybrid Sliding Window) ---
type MovingWindow struct {
	store     Store
	rate      int
	interval  time.Duration
	keyPrefix string
}

func NewMovingWindow(store Store, rate int, intervalSeconds int) *MovingWindow {
	return &MovingWindow{
		store:     store,
		rate:      rate,
		interval:  time.Duration(intervalSeconds) * time.Second,
		keyPrefix: "rate:limiter:moving:",
	}
}

func NewMovingWindowFromConfig(store Store) *MovingWindow {
	return NewMovingWindow(store, cfg.Rate, cfg.Interval)
}

func (mw *MovingWindow) Allow(key string) bool {
	redisKey := fmt.Sprintf("%s%s:%d", mw.keyPrefix, key, time.Now().Unix()/int64(mw.interval.Seconds()))
	count, err := mw.store.Incr(redisKey)
	if err != nil {
		log.Println("Redis INCR error:", err)
		return true
	}
	if count == 1 {
		_ = mw.store.Expire(redisKey, mw.interval*2)
	}
	fmt.Printf("[Moving] Key: %s, Count: %d\n", redisKey, count)
	return int(count) <= mw.rate
}
