package rate_limiter

import "time"

// RedisClientInterface defines the interface for Redis client operations needed by rate limiters
type RedisClientInterface interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	Expire(key string, duration time.Duration, expiryMode ExpiryMode) error
	IncrWithExpiry(key string, duration time.Duration, expiryMode ExpiryMode) (int64, error)
	GetCountAndLastRefill(keyCount, keyLastRefill string) (int64, int, error)
	SetCountAndLastRefill(keyCount, keyLastRefill string, tokenCount int, currentTime int64) error
	HGetAll(key string) (map[string]string, error)
	HIncrByWithExpiry(key string, value string, increment int64, duration time.Duration, expiryMode ExpiryMode) (int64, error)
	HLen(key string) (int64, error)
	HSetWithExpiry(key string, value string, duration time.Duration, expiryMode ExpiryMode) (int64, error)
}
