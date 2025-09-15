package sliding_window_counter_rate_limiter

import (
	"strconv"
	"time"

	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
)

type SlidingWindowCounterRateLimiter struct {
    redisClient rate_limiter.RedisClientInterface
	limit int
	windowSize int64
	subWindowSize int64
}

func NewSlidingWindowCounterRateLimiter(redisClient rate_limiter.RedisClientInterface, limit int, windowSize int64, subWindowSize int64) *SlidingWindowCounterRateLimiter {
	return &SlidingWindowCounterRateLimiter{
		redisClient: redisClient,
		limit: limit,
		windowSize: windowSize,
		subWindowSize: subWindowSize,
	}
}

func (s *SlidingWindowCounterRateLimiter) LimitRequests(clientId string) bool {
	key := "rate_limit:" + clientId
	subWindowCounts, err := s.redisClient.HGetAll(key)
	if err != nil {
		return false
	}

	var totalCount int64
	for _, count := range subWindowCounts {
		c, err := strconv.Atoi(count)
		if err != nil {
			return false
		}

		totalCount += int64(c)
	}

	isAllowed := totalCount < int64(s.limit)
	if isAllowed {
		currentTime := time.Now().Unix()
		currentSubWindow := currentTime / s.subWindowSize

		incrementResult, err := s.redisClient.HIncrByWithExpiry(key, strconv.FormatInt(currentSubWindow, 10), 1, time.Duration(s.subWindowSize)*time.Second, rate_limiter.EXPIRY_MODE_NX)
		if err != nil || incrementResult == 0 {
			return false
		}
	}

	return isAllowed
}