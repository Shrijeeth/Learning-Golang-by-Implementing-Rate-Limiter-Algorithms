package sliding_window_log_rate_limiter

import (
	"time"

	"github.com/google/uuid"

	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
)


type SlidingWindowLogRateLimiter struct {
	redisClient rate_limiter.RedisClientInterface
	limit int
	windowSize int64
}

func NewSlidingWindowLogRateLimiter(redisClient rate_limiter.RedisClientInterface, limit int, windowSize int64) *SlidingWindowLogRateLimiter {
	return &SlidingWindowLogRateLimiter{
		redisClient: redisClient,
		limit: limit,
		windowSize: windowSize,
	}
}

func (s *SlidingWindowLogRateLimiter) LimitRequests(clientId string) bool {
	key := "rate_limit:" + clientId
	fieldKey := uuid.NewString()

	requestCount, err := s.redisClient.HLen(key)
	if err != nil {
		return false
	}

	isAllowed := requestCount < int64(s.limit)
	if isAllowed {
		_, err := s.redisClient.HSetWithExpiry(
			key,
			fieldKey,
			time.Duration(s.windowSize)*time.Second,
			rate_limiter.EXPIRY_MODE_NX,
		)
		if err != nil {
			return false
		}
	}

	return isAllowed
}