package token_bucket_ratelimiter

import (
	"time"

	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
)

type TokenBucketRateLimiter struct {
	redisClient rate_limiter.RedisClientInterface
	bucketCapacity int
	refillRate float64
}

func NewTokenBucketRateLimiter(redisClient rate_limiter.RedisClientInterface, bucketCapacity int, refillRate float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		redisClient: redisClient,
		bucketCapacity: bucketCapacity,
		refillRate: refillRate,
	}
}

func (t *TokenBucketRateLimiter) LimitRequests(clientId string) bool {
	keyCount := "rate_limit:" + clientId + ":count"
	keyLastRefill := "rate_limit:" + clientId + ":lastRefill"
	currentTime := time.Now().Unix()

	lastRefillTime, tokenCount, err := t.redisClient.GetCountAndLastRefill(keyCount, keyLastRefill)
	if err != nil {
		return false
	}

	if lastRefillTime == 0 {
		lastRefillTime = currentTime
		tokenCount = t.bucketCapacity
	}

	elapsedTimeSecs := currentTime - lastRefillTime
	
	tokensToAdd := int(elapsedTimeSecs) * int(t.refillRate)
	tokenCount = min(t.bucketCapacity, tokenCount + tokensToAdd)

	isAllowed := tokenCount > 0
	if isAllowed {
		tokenCount--
	}

	if err := t.redisClient.SetCountAndLastRefill(keyCount, keyLastRefill, tokenCount, currentTime); err != nil {
		return false
	}

	return isAllowed
}