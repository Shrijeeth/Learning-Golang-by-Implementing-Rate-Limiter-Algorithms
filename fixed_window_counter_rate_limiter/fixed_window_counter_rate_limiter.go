package fixed_window_counter_ratelimiter

import (
	"strconv"
	"time"

	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
)

type FixedWindowCounterRateLimiter struct {
    redisClient rate_limiter.RedisClientInterface
	windowSize int
	limit int
}

func NewFixedWindowCounterRateLimiter(redisClient rate_limiter.RedisClientInterface, windowSize int, limit int) *FixedWindowCounterRateLimiter {
	return &FixedWindowCounterRateLimiter{
		redisClient: redisClient,
		windowSize: windowSize,
		limit: limit,
	}
}

func (f *FixedWindowCounterRateLimiter) LimitRequests(clientId string) bool {
	key := "rate_limit:" + clientId
	currentCounterStr, err := f.redisClient.Get(key)
	if err != nil {
		return false
	}

	currentCounter, _ := strconv.Atoi(currentCounterStr)
	if currentCounter >= f.limit {
		return false
	}

	isAllowed := currentCounter < f.limit
	if isAllowed {
		incrResult, err := f.redisClient.IncrWithExpiry(key, time.Duration(f.windowSize) * time.Second, rate_limiter.EXPIRY_MODE_NX)
		if err != nil || incrResult == 0 {
			return false
		}
	}
	return isAllowed
}