package auth

import (
	"sync"
	"time"
)

type TokenBlacklist struct {
	tokens sync.Map
}

var (
	blacklist *TokenBlacklist
	once      sync.Once
)

// GetBlacklist returns the singleton instance of TokenBlacklist
func GetBlacklist() *TokenBlacklist {
	once.Do(func() {
		blacklist = &TokenBlacklist{}
	})
	return blacklist
}

// Add adds a token to the blacklist with an expiration time
func (b *TokenBlacklist) Add(token string, expiry int64) {
	b.tokens.Store(token, expiry)
}

// IsBlacklisted checks if a token is in the blacklist and not expired
func (b *TokenBlacklist) IsBlacklisted(token string) bool {
	if value, exists := b.tokens.Load(token); exists {
		expiry := value.(int64)
		return time.Now().Unix() < expiry
	}
	return false
}

// Add periodic cleanup of expired tokens
func (b *TokenBlacklist) Cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			b.tokens.Range(func(key, value interface{}) bool {
				expiry := value.(int64)
				if time.Now().Unix() > expiry {
					b.tokens.Delete(key)
				}
				return true
			})
		}
	}()
}
