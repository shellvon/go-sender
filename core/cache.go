package core

import (
	"sync"
	"time"
)

// Cache 通用缓存接口.
type Cache[T any] interface {
	// Get 获取缓存值
	Get(key string) (T, bool)
	// Set 存储缓存值，ttl为nil表示永久缓存
	Set(key string, value T, ttl *time.Duration) error
	// Delete 删除缓存项
	Delete(key string) error
}

// CacheItem 缓存项.
type cacheItem struct {
	Value     interface{}
	ExpiresAt *time.Time // nil表示永不过期
}

// IsExpired 检查是否过期.
func (item *cacheItem) IsExpired() bool {
	if item.ExpiresAt == nil {
		return false // 永不过期
	}
	return time.Now().After(*item.ExpiresAt)
}

// MemoryCache 内存缓存实现.
type MemoryCache[T any] struct {
	items map[string]*cacheItem
	mutex sync.RWMutex
}

// NewMemoryCache 创建内存缓存（不启动清理协程）.
func NewMemoryCache[T any]() *MemoryCache[T] {
	return &MemoryCache[T]{
		items: make(map[string]*cacheItem),
	}
}

// NewMemoryCacheWithCleanup 创建内存缓存并启动清理协程.
func NewMemoryCacheWithCleanup[T any](cleanupInterval time.Duration) *MemoryCache[T] {
	cache := &MemoryCache[T]{
		items: make(map[string]*cacheItem),
	}

	// 只有在需要时才启动清理协程
	if cleanupInterval > 0 {
		go cache.cleanupExpired(cleanupInterval)
	}

	return cache
}

// Get 获取缓存值.
func (c *MemoryCache[T]) Get(key string) (T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var zero T
	item, exists := c.items[key]
	if !exists {
		return zero, false
	}

	if item.IsExpired() {
		// 异步删除过期项
		go func() {
			c.mutex.Lock()
			delete(c.items, key)
			c.mutex.Unlock()
		}()
		return zero, false
	}

	if value, ok := item.Value.(T); ok {
		return value, true
	}
	return zero, false
}

// Set 存储缓存值.
func (c *MemoryCache[T]) Set(key string, value T, ttl *time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item := &cacheItem{
		Value:     value,
		ExpiresAt: nil, // 默认不过期
	}

	if ttl != nil {
		expiresAt := time.Now().Add(*ttl)
		item.ExpiresAt = &expiresAt
	}

	c.items[key] = item
	return nil
}

// Delete 删除缓存.
func (c *MemoryCache[T]) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
	return nil
}

// Clear 清空缓存.
func (c *MemoryCache[T]) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}

// Size 获取缓存项数量.
func (c *MemoryCache[T]) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.items)
}

// cleanupExpired 定期清理过期项.
func (c *MemoryCache[T]) cleanupExpired(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, item := range c.items {
			if item.IsExpired() {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}
