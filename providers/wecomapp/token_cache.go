package wecomapp

import (
	"sync"
	"time"
)

// TokenCache 访问令牌存储的接口.
type TokenCache interface {
	Set(key string, token *AccessToken) error
	Get(key string) (*AccessToken, error)
	Delete(key string) error
}

// AccessToken 代表带有过期时间的企业微信访问令牌.
type AccessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired 检查令牌是否已过期或将在缓冲时间内过期.
func (t *AccessToken) IsExpired(buffer time.Duration) bool {
	return time.Now().Add(buffer).After(t.ExpiresAt)
}

// MemoryTokenCache 提供TokenCache的内存实现.
type MemoryTokenCache struct {
	cache map[string]*AccessToken
	mutex sync.RWMutex
}

// NewMemoryTokenCache 创建新的内存令牌缓存.
func NewMemoryTokenCache() *MemoryTokenCache {
	return &MemoryTokenCache{
		cache: make(map[string]*AccessToken),
	}
}

// Set 在缓存中存储令牌.
func (c *MemoryTokenCache) Set(key string, token *AccessToken) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = token
	return nil
}

// Get 从缓存中检索令牌.
func (c *MemoryTokenCache) Get(key string) (*AccessToken, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	token, exists := c.cache[key]
	if !exists {
		return nil, nil
	}
	return token, nil
}

// Delete 从缓存中删除令牌.
func (c *MemoryTokenCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, key)
	return nil
}

// TokenManager handles access token lifecycle management.
type TokenManager struct {
	cache  TokenCache
	buffer time.Duration // Buffer time before token expiration to refresh
}

// NewTokenManager 使用指定的缓存和缓冲时间创建新的token管理器.
func NewTokenManager(cache TokenCache, buffer time.Duration) *TokenManager {
	if cache == nil {
		cache = NewMemoryTokenCache()
	}
	if buffer == 0 {
		buffer = 5 * time.Minute // Default 5 minute buffer
	}
	return &TokenManager{
		cache:  cache,
		buffer: buffer,
	}
}

// GetValidToken 从缓存中检索有效令牌，如果过期/缺失则返回nil.
func (tm *TokenManager) GetValidToken(key string) (*AccessToken, error) {
	token, err := tm.cache.Get(key)
	if err != nil {
		return nil, err
	}

	if token == nil || token.IsExpired(tm.buffer) {
		return nil, nil
	}

	return token, nil
}

// StoreToken stores a new token in the cache.
func (tm *TokenManager) StoreToken(key string, token *AccessToken) error {
	return tm.cache.Set(key, token)
}

// InvalidateToken removes a token from the cache.
func (tm *TokenManager) InvalidateToken(key string) error {
	return tm.cache.Delete(key)
}

// IsTokenValid 检查令牌是否存在且未过期.
func (tm *TokenManager) IsTokenValid(key string) (bool, error) {
	token, err := tm.GetValidToken(key)
	if err != nil {
		return false, err
	}
	return token != nil, nil
}
