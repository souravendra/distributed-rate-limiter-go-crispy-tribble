package limiter_test

import (
	"testing"
	"time"

	"github.com/souravendra/distributed-rate-limiter-go-crispy-tribble/limiter"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	counter map[string]int64
	ttl     map[string]time.Duration
}

func newMockStore() *mockStore {
	return &mockStore{
		counter: make(map[string]int64),
		ttl:     make(map[string]time.Duration),
	}
}

func (m *mockStore) Incr(key string) (int64, error) {
	m.counter[key]++
	return m.counter[key], nil
}

func (m *mockStore) Expire(key string, ttl time.Duration) error {
	m.ttl[key] = ttl
	return nil
}

func (m *mockStore) GetTTL(key string) (time.Duration, error) {
	return m.ttl[key], nil
}

func TestTokenBucket_Allow(t *testing.T) {
	store := newMockStore()
	bucket := limiter.NewTokenBucket(store, 2, 1, 2)
	key := "test"

	allowed := bucket.Allow(key)
	assert.True(t, allowed)

	allowed = bucket.Allow(key)
	assert.True(t, allowed)

	allowed = bucket.Allow(key)
	assert.False(t, allowed)
}
