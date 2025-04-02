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
	t.Parallel()
	store := newMockStore()
	tb := limiter.NewTokenBucket(store, 2, 1, 2)
	key := "token"

	assert.True(t, tb.Allow(key))
	assert.True(t, tb.Allow(key))
	assert.False(t, tb.Allow(key))
}

func TestFixedWindow_Allow(t *testing.T) {
	t.Parallel()
	store := newMockStore()
	fw := limiter.NewFixedWindow(store, 2, 1)
	key := "fixed"

	assert.True(t, fw.Allow(key))
	assert.True(t, fw.Allow(key))
	assert.False(t, fw.Allow(key))
}

func TestSlidingWindowLog_Allow(t *testing.T) {
	t.Parallel()
	store := newMockStore()
	sw := limiter.NewSlidingWindowLog(store, 2, 1)
	key := "sliding"

	assert.True(t, sw.Allow(key))
	assert.True(t, sw.Allow(key))
	assert.False(t, sw.Allow(key))
}

func TestLeakyBucket_Allow(t *testing.T) {
	t.Parallel()
	store := newMockStore()
	lb := limiter.NewLeakyBucket(store, 2, 1, 2)
	key := "leaky"

	assert.True(t, lb.Allow(key))
	assert.True(t, lb.Allow(key))
	assert.False(t, lb.Allow(key))
}

func TestMovingWindow_Allow(t *testing.T) {
	t.Parallel()
	store := newMockStore()
	mw := limiter.NewMovingWindow(store, 2, 1)
	key := "moving"

	assert.True(t, mw.Allow(key))
	assert.True(t, mw.Allow(key))
	assert.False(t, mw.Allow(key))
}
