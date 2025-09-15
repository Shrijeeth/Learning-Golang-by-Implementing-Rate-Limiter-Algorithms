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
	
	// If there's an error and it's not just an empty string (new client), reject the request
	if err != nil {
		return false
	}

	// For new clients or expired windows, currentCounterStr will be empty
	currentCounter := 0
	if currentCounterStr != "" {
		currentCounter, _ = strconv.Atoi(currentCounterStr)
	}

	// If counter is at or above limit, reject the request
	if currentCounter >= f.limit {
		return false
	}

	// Request is allowed, increment the counter and set expiry
	incrResult, err := f.redisClient.IncrWithExpiry(key, time.Duration(f.windowSize) * time.Second, rate_limiter.EXPIRY_MODE_NX)
	if err != nil || incrResult == 0 {
		return false
	}

	return true
}