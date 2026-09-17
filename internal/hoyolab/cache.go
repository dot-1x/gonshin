package hoyolab

import (
	"sync"
	"time"
)

type cacheEntry struct {
	value   any
	expires time.Time
}

type cacheCall struct {
	done  chan struct{}
	value any
	err   error
}

// ttlCache is a small in-memory cache with per-key request coalescing. It
// mirrors the original Next.js `revalidate` behaviour (successful responses
// cached for a TTL, transient failures not cached).
type ttlCache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
	calls   map[string]*cacheCall
}

func newTTLCache() *ttlCache {
	return &ttlCache{
		entries: make(map[string]cacheEntry),
		calls:   make(map[string]*cacheCall),
	}
}

func (c *ttlCache) get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expires) {
		delete(c.entries, key)
		return nil, false
	}
	return e.value, true
}

func (c *ttlCache) set(key string, value any, ttl time.Duration) {
	if ttl <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{value: value, expires: time.Now().Add(ttl)}
}

// do returns the cached value for key, or runs fetch (deduplicating concurrent
// callers) and caches the result on success.
func (c *ttlCache) do(key string, ttl time.Duration, fetch func() (any, error)) (any, error) {
	if v, ok := c.get(key); ok {
		return v, nil
	}

	c.mu.Lock()
	if call, ok := c.calls[key]; ok {
		c.mu.Unlock()
		<-call.done
		return call.value, call.err
	}
	call := &cacheCall{done: make(chan struct{})}
	c.calls[key] = call
	c.mu.Unlock()

	call.value, call.err = fetch()
	if call.err == nil && call.value != nil {
		c.set(key, call.value, ttl)
	}

	c.mu.Lock()
	delete(c.calls, key)
	c.mu.Unlock()
	close(call.done)

	return call.value, call.err
}
