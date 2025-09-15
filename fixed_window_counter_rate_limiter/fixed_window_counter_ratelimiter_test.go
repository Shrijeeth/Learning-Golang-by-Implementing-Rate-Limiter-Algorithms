package fixed_window_counter_ratelimiter_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixed_window_counter_ratelimiter "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/fixed_window_counter_rate_limiter"
	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter/mocks"
)

var _ = Describe("FixedWindowCounterRatelimiter", func() {
	var (
		mockRedisClient *mocks.MockRedisClient
		rateLimiter     *fixed_window_counter_ratelimiter.FixedWindowCounterRateLimiter
		clientID        string
		windowSize      int
		limit           int
	)

	BeforeEach(func() {
		mockRedisClient = mocks.NewMockRedisClient()
		clientID = "test-client"
		windowSize = 10
		limit = 5
	})

	Describe("LimitRequests", func() {
		Context("when client is new (first request)", func() {
			It("should allow the request and set expiry", func() {
				// Mock the Get function to simulate a new client (no existing counter)
				mockRedisClient.GetFunc = func(key string) (string, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return "", nil // Empty string for a new client
				}

				// Mock the IncrWithExpiry function to simulate a successful increment
				mockRedisClient.IncrWithExpiryFunc = func(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					Expect(duration).To(Equal(time.Duration(windowSize) * time.Second))
					Expect(expiryMode).To(Equal(rate_limiter.EXPIRY_MODE_NX))
					return 1, nil // First request, counter is 1
				}

				rateLimiter = fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(mockRedisClient, windowSize, limit)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
			})
		})

		Context("when client is within the limit", func() {
			It("should allow the request if counter is less than limit", func() {
				// Mock the Get function to simulate an existing counter
				mockRedisClient.GetFunc = func(key string) (string, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return "3", nil // Current counter is 3, which is below limit of 5
				}

				// Mock the IncrWithExpiry function
				mockRedisClient.IncrWithExpiryFunc = func(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return 4, nil // Increment to 4
				}

				rateLimiter = fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(mockRedisClient, windowSize, limit)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
			})
		})

		Context("when client reaches the limit", func() {
			It("should allow the request if counter equals limit", func() {
				// Mock the Get function to simulate an existing counter
				mockRedisClient.GetFunc = func(key string) (string, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return "4", nil // Current counter is 4, which is below limit of 5
				}

				// Mock the IncrWithExpiry function
				mockRedisClient.IncrWithExpiryFunc = func(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return 5, nil // Increment to 5 (equal to limit)
				}

				rateLimiter = fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(mockRedisClient, windowSize, limit)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
			})
		})

		Context("when client exceeds the limit", func() {
			It("should reject the request", func() {
				// Mock the Get function to simulate an existing counter
				mockRedisClient.GetFunc = func(key string) (string, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return "5", nil // Current counter is 5, which equals the limit
				}

				rateLimiter = fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(mockRedisClient, windowSize, limit)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeFalse())
			})
		})

		Context("when Redis returns an error on Get", func() {
			It("should reject the request", func() {
				// Mock the Get function to simulate a Redis error
				mockRedisClient.GetFunc = func(key string) (string, error) {
					return "", errors.New("redis connection error")
				}

				rateLimiter = fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(mockRedisClient, windowSize, limit)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeFalse())
			})
		})

		Context("when Redis returns an error on IncrWithExpiry", func() {
			It("should reject the request", func() {
				// Mock the Get function to simulate a valid counter
				mockRedisClient.GetFunc = func(key string) (string, error) {
					return "3", nil
				}

				// Mock the IncrWithExpiry function to simulate a Redis error
				mockRedisClient.IncrWithExpiryFunc = func(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
					return 0, errors.New("redis write error")
				}

				rateLimiter = fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(mockRedisClient, windowSize, limit)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeFalse())
			})
		})

		Context("when window expires and resets", func() {
			It("should allow requests in a new window", func() {
				// First call - simulate that the key doesn't exist (window expired)
				mockRedisClient.GetFunc = func(key string) (string, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return "", nil // Empty string but no error, which is handled as a new window
				}

				// Mock the IncrWithExpiry function for a new window
				mockRedisClient.IncrWithExpiryFunc = func(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
					Expect(key).To(Equal("rate_limit:test-client"))
					return 1, nil // First request in new window
				}

				rateLimiter = fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(mockRedisClient, windowSize, limit)
				result := rateLimiter.LimitRequests(clientID)

				Expect(result).To(BeTrue())
			})
		})
	})
})
