package token_bucket_ratelimiter_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	token_bucket_ratelimiter "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/token_bucket_rate_limiter"
	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter/mocks"
)

var _ = Describe("TokenBucketRatelimiter", func() {
	var (
		mockRedisClient *mocks.MockRedisClient
		rateLimiter     *token_bucket_ratelimiter.TokenBucketRateLimiter
		clientID        string
		bucketCapacity  int
		refillRate      float64
		currentTime     int64
	)

	BeforeEach(func() {
		mockRedisClient = mocks.NewMockRedisClient()
		clientID = "test-client"
		bucketCapacity = 10
		refillRate = 1.0 // 1 token per second
		currentTime = time.Now().Unix()
	})

	Describe("LimitRequests", func() {
		Context("when client is new (first request)", func() {
			It("should fill the bucket to capacity and allow the request", func() {
				// Mock Redis client to return empty data (new client)
				mockRedisClient.GetCountAndLastRefillFunc = func(keyCount, keyLastRefill string) (int64, int, error) {
					return 0, 0, nil
				}

				var capturedTokenCount int
				var capturedTime int64

				mockRedisClient.SetCountAndLastRefillFunc = func(keyCount, keyLastRefill string, tokenCount int, time int64) error {
					capturedTokenCount = tokenCount
					capturedTime = time
					return nil
				}

				rateLimiter = token_bucket_ratelimiter.NewTokenBucketRateLimiter(mockRedisClient, bucketCapacity, refillRate)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
				Expect(capturedTokenCount).To(Equal(bucketCapacity - 1)) // Initial capacity minus one token for the request
				Expect(capturedTime).To(BeNumerically(">", 0))
			})
		})

		Context("when tokens are available", func() {
			It("should allow the request and decrement token count", func() {
				initialTokens := 5

				mockRedisClient.GetCountAndLastRefillFunc = func(keyCount, keyLastRefill string) (int64, int, error) {
					return currentTime - 10, initialTokens, nil
				}

				var capturedTokenCount int
				mockRedisClient.SetCountAndLastRefillFunc = func(keyCount, keyLastRefill string, tokenCount int, time int64) error {
					capturedTokenCount = tokenCount
					return nil
				}

				rateLimiter = token_bucket_ratelimiter.NewTokenBucketRateLimiter(mockRedisClient, bucketCapacity, refillRate)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
				// The calculation is: min(bucketCapacity, initialTokens + tokensToAdd) - 1
				// initialTokens = 5, tokensToAdd = 10 (elapsed time * refill rate), but capped at bucketCapacity = 10
				// So it's min(10, 5 + 10) - 1 = min(10, 15) - 1 = 10 - 1 = 9
				Expect(capturedTokenCount).To(Equal(9))
			})
		})

		Context("when tokens are refilled", func() {
			It("should refill tokens based on elapsed time", func() {
				elapsedSeconds := 5
				initialTokens := 2

				mockRedisClient.GetCountAndLastRefillFunc = func(keyCount, keyLastRefill string) (int64, int, error) {
					return currentTime - int64(elapsedSeconds), initialTokens, nil
				}

				var capturedTokenCount int
				mockRedisClient.SetCountAndLastRefillFunc = func(keyCount, keyLastRefill string, tokenCount int, time int64) error {
					capturedTokenCount = tokenCount
					return nil
				}

				rateLimiter = token_bucket_ratelimiter.NewTokenBucketRateLimiter(mockRedisClient, bucketCapacity, refillRate)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
				// The calculation is: min(bucketCapacity, initialTokens + tokensToAdd) - 1
				// initialTokens = 2, tokensToAdd = 5 (elapsed time * refill rate), bucketCapacity = 10
				// So it's min(10, 2 + 5) - 1 = min(10, 7) - 1 = 7 - 1 = 6
				Expect(capturedTokenCount).To(Equal(6))
			})

			It("should not exceed bucket capacity when refilling", func() {
				elapsedSeconds := 20 // More than enough to fill the bucket
				initialTokens := 2

				mockRedisClient.GetCountAndLastRefillFunc = func(keyCount, keyLastRefill string) (int64, int, error) {
					return currentTime - int64(elapsedSeconds), initialTokens, nil
				}

				var capturedTokenCount int
				mockRedisClient.SetCountAndLastRefillFunc = func(keyCount, keyLastRefill string, tokenCount int, time int64) error {
					capturedTokenCount = tokenCount
					return nil
				}

				rateLimiter = token_bucket_ratelimiter.NewTokenBucketRateLimiter(mockRedisClient, bucketCapacity, refillRate)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
				Expect(capturedTokenCount).To(Equal(bucketCapacity - 1)) // Full bucket minus one consumed token
			})
		})

		Context("when bucket is empty", func() {
			It("should reject the request", func() {
				mockRedisClient.GetCountAndLastRefillFunc = func(keyCount, keyLastRefill string) (int64, int, error) {
					return currentTime, 0, nil // No tokens available
				}

				var capturedTokenCount int
				mockRedisClient.SetCountAndLastRefillFunc = func(keyCount, keyLastRefill string, tokenCount int, time int64) error {
					capturedTokenCount = tokenCount
					return nil
				}

				rateLimiter = token_bucket_ratelimiter.NewTokenBucketRateLimiter(mockRedisClient, bucketCapacity, refillRate)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeFalse())
				Expect(capturedTokenCount).To(Equal(0)) // Should remain empty
			})
		})

		Context("when Redis client returns an error", func() {
			It("should reject the request on GetCountAndLastRefill error", func() {
				mockRedisClient.GetCountAndLastRefillFunc = func(keyCount, keyLastRefill string) (int64, int, error) {
					return 0, 0, errors.New("redis connection error")
				}

				rateLimiter = token_bucket_ratelimiter.NewTokenBucketRateLimiter(mockRedisClient, bucketCapacity, refillRate)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeFalse())
			})

			It("should reject the request on SetCountAndLastRefill error", func() {
				mockRedisClient.GetCountAndLastRefillFunc = func(keyCount, keyLastRefill string) (int64, int, error) {
					return currentTime - 10, 5, nil
				}

				mockRedisClient.SetCountAndLastRefillFunc = func(keyCount, keyLastRefill string, tokenCount int, time int64) error {
					return errors.New("redis write error")
				}

				rateLimiter = token_bucket_ratelimiter.NewTokenBucketRateLimiter(mockRedisClient, bucketCapacity, refillRate)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeFalse())
			})
		})
	})
})
