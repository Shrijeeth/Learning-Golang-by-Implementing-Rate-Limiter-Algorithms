package rate_limiter

// RedisClientInterface defines the interface for Redis client operations needed by rate limiters
type RedisClientInterface interface {
	GetCountAndLastRefill(keyCount, keyLastRefill string) (int64, int, error)
	SetCountAndLastRefill(keyCount, keyLastRefill string, tokenCount int, currentTime int64) error
}
