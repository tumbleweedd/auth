package inmemory

import (
	"context"
	"sync"
	"time"
)

type item[V any] struct {
	value  V
	expiry time.Time
}

func (i *item[V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

type Cache[K comparable, V any] struct {
	items      map[K]*item[V]
	mu         sync.RWMutex
	defaultTTL time.Duration
}

func New[K comparable, V any](ctx context.Context, defaultTTL time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		items:      make(map[K]*item[V]),
		defaultTTL: defaultTTL,
	}

	go func() {
		timeDuration := 5 * time.Minute
		timer := time.NewTimer(timeDuration)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				cache.mu.Lock()
				for key, cacheItem := range cache.items {
					if cacheItem.isExpired() {
						delete(cache.items, key)
					}
				}
				cache.mu.Unlock()
			}
		}
	}()

	return cache
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.RLock()

	cacheItem, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return cacheItem.value, false
	}

	if cacheItem.isExpired() {
		c.mu.RUnlock()
		return cacheItem.value, false
	}

	c.mu.RUnlock()

	return cacheItem.value, true
}

func (c *Cache[K, V]) Add(key K, value V) {
	c.mu.Lock()

	c.items[key] = &item[V]{
		value:  value,
		expiry: time.Now().Add(c.defaultTTL),
	}

	c.mu.Unlock()
}

func (c *Cache[K, V]) Range(fn func(key K, value V) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for key, cacheItem := range c.items {
		if !fn(key, cacheItem.value) {
			break
		}
	}
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}
